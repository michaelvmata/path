package world

import (
	"errors"
	"fmt"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/memory"
	"github.com/michaelvmata/path/modifiers"
	"github.com/michaelvmata/path/session"
	"github.com/michaelvmata/path/skills"
	"github.com/michaelvmata/path/stats"
	"github.com/michaelvmata/path/symbols"
	"log"
	"strings"
)

type Buff interface {
	Update(int)
	Expire()
	IsExpired() bool
	Name() string
	ApplyMessage() string
	UnapplyMessage() string
	AlreadyApplied() string
	Upkeep() int
	Remaining() int
}

type CoolDown interface {
	Update(int)
	IsExpired() bool
	Name() string
}

type Character struct {
	UUID    string
	Name    string
	Room    *Room
	Session *session.Session

	keywords []string

	Health stats.Line
	Spirit stats.Line

	Core   stats.Core
	Skills skills.Skills

	Gear      *item.Gear
	Inventory item.Container

	Essence int

	Attacking []*Character

	Buffs     []Buff
	CoolDowns []CoolDown
	Memory    *memory.Memory

	IsAggressive bool
	IsSocial     bool
	IsPlayer     bool

	Stunned int
}

func (c *Character) CreditEssence(amount int) {
	c.Essence += amount
}

func (c *Character) DebitEssence(amount int) {
	c.Essence -= amount
}

func (c *Character) Apply(buff Buff) {
	for _, b := range c.Buffs {
		if b.Name() == buff.Name() {
			c.Showln(buff.AlreadyApplied())
			return
		}
	}
	c.Buffs = append(c.Buffs, buff)
	c.Showln(buff.ApplyMessage())
}

func (c *Character) Unapply(buffName string) {
	remaining := make([]Buff, 0)
	for _, buff := range c.Buffs {
		if buff.Name() != buffName {
			remaining = append(remaining, buff)
		} else {
			c.Showln(buff.UnapplyMessage())
		}
	}
	if len(remaining) == len(c.Buffs) {
		log.Fatalf("Tried to unapply buff that wasn't applied %s", buffName)
	}
	c.Buffs = remaining
}

func (c *Character) HasBuff(buffName string) bool {
	for _, buff := range c.Buffs {
		if buff.Name() == buffName {
			return true
		}
	}
	return false
}

func (c *Character) UnapplyExpiredBuffs() {
	buffs := make([]Buff, 0)
	messages := make([]string, 0)
	for _, b := range c.Buffs {
		if !b.IsExpired() {
			buffs = append(buffs, b)
		} else {
			messages = append(messages, b.UnapplyMessage())
			log.Printf("Expiring %s from %s", b.Name(), c.Name)
		}
	}
	c.Buffs = buffs
	if len(messages) > 0 {
		c.ShowNewline()
		message := strings.Join(messages, "\n")
		c.Showln(message)
		c.ShowNewline()
		c.ShowPrompt()
	}
}

func (c *Character) ApplyCoolDown(coolDown CoolDown) {
	for _, cd := range c.CoolDowns {
		if cd.Name() == coolDown.Name() {
			c.Showln("CoolDown already applied")
			return
		}
	}
	c.CoolDowns = append(c.CoolDowns, coolDown)
}

func (c *Character) UnapplyExpiredCoolDowns() {
	coolDowns := make([]CoolDown, 0)
	for _, cd := range c.CoolDowns {
		if !cd.IsExpired() {
			coolDowns = append(coolDowns, cd)
		} else {
			c.Showln("Removing cool down %s", cd.Name())
		}
	}
	c.CoolDowns = coolDowns
}

func (c *Character) OnCoolDown(name string) bool {
	for _, cd := range c.CoolDowns {
		if cd.Name() == name {
			return true
		}
	}
	return false
}

func (c *Character) Stun(length int) {
	c.Stunned += length
}

func (c *Character) ReduceStun() {
	if c.Stunned > 0 {
		c.Stunned--
	}
}

func (c *Character) Aggro() {
	if c.IsFighting() || c.IsPlayer || !c.IsAggressive {
		return
	}
	for _, candidate := range c.Room.Players {
		if !candidate.IsPlayer {
			continue
		}
		c.StartAttacking(candidate)
		candidate.StartAttacking(c)

		message := Message{
			FirstPerson:        c,
			FirstPersonMessage: "You scream, \"This is SPARTA!\"",
			ThirdPersonMessage: fmt.Sprintf("%s screams, \"This is SPARTA!\"", c.Name),
		}
		c.Room.ShowMessage(message)
		break
	}
}

func (c *Character) Social() {
	if c.IsFighting() || c.IsPlayer || !c.IsSocial {
		return
	}
	for _, candidate := range c.Room.Players {
		if candidate.IsPlayer || !candidate.IsFighting() {
			continue
		}
		target := candidate.ImmediateDefender()
		if target == nil {
			log.Fatalf("Mobile %s is fighting without a defender", candidate.Name)
		}
		c.StartAttacking(target)
		target.StartAttacking(c)

		message := Message{
			FirstPerson:        c,
			FirstPersonMessage: "You scream, \"And my AXE!\"",
			ThirdPersonMessage: fmt.Sprintf("%s bellows, \"And my AXE!\"", c.Name),
		}
		c.Room.ShowMessage(message)

		// bail now that character is fighting
		break
	}
}

func (c *Character) IsStunned() bool {
	return c.Stunned > 0
}

func (c *Character) Weapon() *item.Weapon {
	if c.Gear.MainHand != nil {
		return c.Gear.MainHand
	}
	hand := item.NewWeapon("Barehand", "Barehand", []string{}, item.Crush, []string{item.Impact})
	hand.MinimumDamage = 1
	hand.MaximumDamage = 5
	hand.CriticalBonus = 1
	hand.CriticalRate = 0
	return hand
}

func (c *Character) Clone(target Character) {
	c.UUID = target.UUID
	c.Name = target.Name
	c.Room = nil
	c.Session = nil
	c.Essence = target.Essence

	c.Core.Power.Base = target.Core.Power.Base
	c.Core.Agility.Base = target.Core.Agility.Base
	c.Core.Insight.Base = target.Core.Insight.Base
	c.Core.Will.Base = target.Core.Will.Base

	c.IsAggressive = target.IsAggressive
	c.IsSocial = target.IsSocial

	c.Gear = item.NewGear()
	c.Inventory = item.NewContainer(10)
	c.Attacking = make([]*Character, 0)
}

func NewPlayer(UUID string, handle string) *Character {
	return &Character{
		UUID: UUID,
		Name: handle,

		keywords: strings.Fields(strings.ToLower(handle)),

		Health: stats.Line{},
		Spirit: stats.Line{},

		Core:   stats.NewCore(),
		Skills: skills.NewSkills(),

		Gear:      item.NewGear(),
		Inventory: item.NewContainer(10),

		Attacking: make([]*Character, 0),

		Memory:       memory.NewMemory(),
		IsAggressive: false,
		IsSocial:     false,
		IsPlayer:     true,
	}
}

func (c *Character) ImmediateDefender() *Character {
	if !c.IsFighting() {
		return nil
	}
	return c.Attacking[0]
}

func (c *Character) ShowPrompt() {
	border := "<grey_62>>> "
	c.Show("%s%s <red>%d%s <green>%d%s %s",
		border,
		c.Name,
		c.Health.Current, symbols.HEART,
		c.Spirit.Current, symbols.TWELVE_STAR,
		border)
	if !c.IsFighting() {
		return
	}
	target := c.ImmediateDefender()
	health := int((float64(target.Health.Current) / float64(target.Health.Maximum)) * 100)
	c.Show("%s <blue>%d%s %s", target.Name, health, symbols.CIRCLED_BULLET, border)
}

func (c *Character) HasKeyword(target string) bool {
	target = strings.ToLower(target)
	for _, keyword := range c.keywords {
		if keyword == target {
			return true
		}
	}
	return false
}

func (c *Character) StartAttacking(defender *Character) {
	for _, target := range c.Attacking {
		if target == defender {
			return
		}
	}
	c.Attacking = append(c.Attacking, defender)
}

func (c *Character) StopAttacking(defender *Character) {
	index := -1
	for i, target := range c.Attacking {
		if target == defender {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	c.Attacking = append(c.Attacking[:index], c.Attacking[index+1:]...)
}

func (c *Character) IsAttacking(defender *Character) bool {
	for _, target := range c.Attacking {
		if target == defender {
			return true
		}
	}
	return false
}

func (c *Character) IsFighting() bool {
	return len(c.Attacking) > 0
}

func (c *Character) Discard(keyword string) item.Item {
	index := c.Inventory.IndexOfItem(keyword)
	if index == -1 {
		return nil
	}
	i := c.Inventory.RemItemAtIndex(index)
	return i
}

func (c *Character) Receive(i item.Item) error {
	if err := c.Inventory.AddItem(i); err != nil {
		return errors.New("player can't carry item")
	}
	return nil
}

func (c *Character) Move(r *Room) {
	c.Room = r
}

func (c *Character) ApplyModifiers(mods []modifiers.Modifier) {
	for _, mod := range mods {
		switch mod.Type {
		case modifiers.Power:
			c.Core.Power.Modify(mod.Value)
		case modifiers.Agility:
			c.Core.Agility.Modify(mod.Value)
		case modifiers.Insight:
			c.Core.Insight.Modify(mod.Value)
		case modifiers.Will:
			c.Core.Will.Modify(mod.Value)
		case modifiers.Dagger:
			c.Skills.Dagger.Modify(mod.Value)
		case modifiers.Sword:
			c.Skills.Sword.Modify(mod.Value)
		case modifiers.Spear:
			c.Skills.Spear.Modify(mod.Value)
		}
	}
}

func (c *Character) ApplyItemModifiers(i item.Item) {
	if !item.IsNil(i) {
		c.ApplyModifiers(i.Modifiers())
	}
}

func (c *Character) CalculateModifiers() {
	c.Core.ResetModifier()
	c.ApplyItemModifiers(c.Gear.Head)
	c.ApplyItemModifiers(c.Gear.Neck)
	c.ApplyItemModifiers(c.Gear.Body)
	c.ApplyItemModifiers(c.Gear.Arms)
	c.ApplyItemModifiers(c.Gear.Hands)
	c.ApplyItemModifiers(c.Gear.Waist)
	c.ApplyItemModifiers(c.Gear.Legs)
	c.ApplyItemModifiers(c.Gear.Feet)
	c.ApplyItemModifiers(c.Gear.Wrist)
	c.ApplyItemModifiers(c.Gear.Fingers)
	c.ApplyItemModifiers(c.Gear.OffHand)
	c.ApplyItemModifiers(c.Gear.MainHand)
}

func (c *Character) IsDead() bool {
	return c.Health.Current <= 0
}

func (c *Character) Restore() {
	c.CalculateModifiers()
	// Adjust lines from core stats

	c.Health.Maximum = c.Core.Power.Value() * 100
	c.Health.Current = c.Health.Maximum
	c.Health.RecoverRate = c.Core.Power.Value()
	c.Spirit.Maximum = c.Core.Will.Value() * 100
	c.Spirit.Current = c.Spirit.Maximum
	c.Spirit.RecoverRate = c.Core.Will.Value()
}

func (c *Character) Update(tick int) {
	c.CalculateModifiers()
	// Adjust lines from core stats

	c.Health.Maximum = c.Core.Power.Value() * 100
	c.Health.EnforceMaximum()
	c.Health.RecoverRate = c.Core.Power.Value()
	c.Spirit.Maximum = c.Core.Will.Value() * 100
	c.Spirit.EnforceMaximum()
	c.Spirit.RecoverRate = c.Core.Will.Value()

	if tick > 0 {
		if tick%5 == 0 && !c.IsFighting() {
			c.Health.Recover()
			c.Spirit.Recover()
		}
		for _, buff := range c.Buffs {
			buff.Update(tick)
		}
		c.UnapplyExpiredBuffs()
		for _, cd := range c.CoolDowns {
			cd.Update(tick)
		}
		c.UnapplyExpiredCoolDowns()
		c.ReduceStun()
		c.Aggro()
		c.Social()

		c.Memory.Update(tick)
	}
}

func (c *Character) ShowNewline() {
	if c.Session != nil {
		c.Session.Outgoing <- "\n"
	}
}

func (c *Character) Showln(message string, args ...interface{}) {
	if c.Session != nil {
		c.Session.Outgoing <- fmt.Sprintf(message, args...)
		c.Session.Outgoing <- "\n"
	}
}

func (c *Character) ShowDivider() {
	if c.Session != nil {
		c.Session.Outgoing <- "------------------------------------------------------------"
	}
}

func (c *Character) Show(message string, args ...interface{}) {
	if c.Session != nil {
		c.Session.Outgoing <- fmt.Sprintf(message, args...)
	}
}

func (c *Character) Describe() string {
	if len(c.Attacking) > 0 {
		return fmt.Sprintf("%s is fighting.", c.Name)
	}
	return fmt.Sprintf("%s is here.", c.Name)
}

type RoomMobile struct {
	MobileUUID string
	Count      int
}

func NewRoomMobile(mobileUUID string, count int) RoomMobile {
	return RoomMobile{
		MobileUUID: mobileUUID,
		Count:      count,
	}
}

type Exits struct {
	East  string
	North string
	South string
	West  string
}

func (e Exits) FirstExit() string {
	if e.East != "" {
		return e.East
	}
	if e.North != "" {
		return e.North
	}
	if e.South != "" {
		return e.South
	}
	if e.West != "" {
		return e.West
	}
	return ""
}

type Room struct {
	UUID        string
	name        string
	description string
	Players     []*Character
	Items       item.Container
	Size        int
	Exits       Exits
}

func NewRoom(uuid string, name string, description string, size int) *Room {
	room := Room{
		UUID:        uuid,
		name:        name,
		description: description,
		Size:        size,
		Players:     make([]*Character, 0, size),
		Items:       item.NewContainer(100),
		Exits:       Exits{},
	}
	return &room
}

func (r *Room) ShowMessage(message Message) error {
	if message.FirstPerson != nil && message.FirstPersonMessage != "" {
		if r.IndexOfPlayer(message.FirstPerson) == -1 {
			return errors.New("first person not in room")
		}
		message.FirstPerson.Showln(message.FirstPersonMessage)
	}
	if message.SecondPerson != nil && message.SecondPersonMessage != "" {
		if r.IndexOfPlayer(message.SecondPerson) == -1 {
			return errors.New("second person not in room")
		}
		message.SecondPerson.Showln(message.SecondPersonMessage)
	}
	if message.ThirdPersonMessage == "" {
		return nil
	}
	for _, player := range r.Players {
		if player == message.FirstPerson || player == message.SecondPerson {
			continue
		}
		player.Showln(message.ThirdPersonMessage)
	}
	return nil
}

func (r *Room) Describe(firstPerson *Character) string {
	parts := make([]string, 0)
	parts = append(parts, r.name)
	parts = append(parts, "")
	parts = append(parts, r.description)
	parts = append(parts, "")
	parts = append(parts, fmt.Sprintf("[%s]", r.DescribeExits()))
	parts = append(parts, "")
	for _, i := range r.Items.Items {
		parts = append(parts, i.Name())
	}
	for _, p := range r.Players {
		if firstPerson != p {
			parts = append(parts, p.Describe())
		}
	}
	return strings.Join(parts, "\n")
}

func (r *Room) DescribeExits() string {
	parts := make([]string, 0)
	if r.Exits.East != "" {
		parts = append(parts, "east")
	}
	if r.Exits.North != "" {
		parts = append(parts, "north")
	}
	if r.Exits.South != "" {
		parts = append(parts, "south")
	}
	if r.Exits.West != "" {
		parts = append(parts, "west")
	}
	if len(parts) == 0 {
		parts = append(parts, "None")
	}
	return strings.Join(parts, ", ")
}

func (r *Room) Accept(i item.Item) error {
	// Accept places an item in the room.  It throws an error if there is no space.
	err := r.Items.AddItem(i)
	if err != nil {
		return errors.New("can't drop item")
	}
	return nil
}

func (r *Room) PickupItem(keyword string) (item.Item, error) {
	index := r.Items.IndexOfItem(keyword)
	if index == -1 {
		return nil, errors.New("no item with keyword")
	}
	i := r.Items.RemItemAtIndex(index)
	return i, nil
}

func (r *Room) IsFull() bool {
	return r.Size == len(r.Players)
}

func (r *Room) Enter(c *Character) error {
	if r.IsFull() {
		return errors.New("room is full")
	}
	if r.IndexOfPlayer(c) != -1 {
		return errors.New("player already in room")
	}
	r.Players = append(r.Players, c)
	return nil
}

func (r *Room) Exit(c *Character) error {
	i := r.IndexOfPlayer(c)
	if i == -1 {
		return errors.New("player not in room")
	}
	r.Players = append(r.Players[:i], r.Players[i+1:]...)
	return nil
}

func (r *Room) IndexOfPlayerHandle(handle string) int {
	for i, p := range r.Players {
		if p.HasKeyword(handle) {
			return i
		}
	}
	return -1
}

func (r *Room) GetPlayer(handle string) *Character {
	index := r.IndexOfPlayerHandle(handle)
	if index == -1 {
		return nil
	}
	return r.Players[index]
}

func (r *Room) IndexOfPlayer(target *Character) int {
	for i, p := range r.Players {
		if p == target {
			return i
		}
	}
	return -1
}

func (r *Room) MobileCount(mobileUUID string) int {
	// Count number of mobiles with UUID in room
	count := 0
	for _, player := range r.Players {
		if player.UUID == mobileUUID {
			count += 1
		}
	}
	return count
}

type Mobiles struct {
	Prototypes map[string]Character
	Instances  []*Character
}

func (m *Mobiles) AddPrototype(p Character) {
	m.Prototypes[p.UUID] = p
}

func (m *Mobiles) Spawn(UUID string) *Character {
	prototype := m.Prototypes[UUID]
	mobile := NewPlayer(UUID, prototype.Name)
	mobile.IsPlayer = false
	mobile.Clone(prototype)
	m.Instances = append(m.Instances, mobile)
	return mobile
}

func (m *Mobiles) Unspawn(c *Character) error {
	index := -1
	for i, m := range m.Instances {
		if m == c {
			index = i
		}
	}
	// Players wil be not found
	if index != -1 {
		m.Instances = append(m.Instances[:index], m.Instances[index+1:]...)
		return nil
	}
	return errors.New("not a mobile")
}

func (m *Mobiles) IsInstance(c *Character) bool {
	for _, instance := range m.Instances {
		if instance == c {
			return true
		}
	}
	return false
}

type World struct {
	Players     map[string]*Character
	Mobiles     Mobiles
	Rooms       map[string]*Room
	RoomMobiles map[string][]RoomMobile
	Items       map[string]item.Item

	Ticks       int
	SpawnTicks  int
	BattleTicks int
}

func NewWorld() *World {
	w := World{
		Players: make(map[string]*Character),
		Mobiles: Mobiles{
			Prototypes: make(map[string]Character),
			Instances:  make([]*Character, 0),
		},
		Rooms:       make(map[string]*Room),
		RoomMobiles: make(map[string][]RoomMobile, 0),
		Items:       make(map[string]item.Item),
		SpawnTicks:  60,
		BattleTicks: 3,
	}
	return &w
}

func (w *World) IsMobile(c *Character) bool {
	return w.Mobiles.IsInstance(c)
}

func (w *World) Update() {
	w.Ticks++
	for _, player := range w.Players {
		player.Update(w.Ticks)
	}
	for _, mobile := range w.Mobiles.Instances {
		mobile.Update(w.Ticks)
	}
	if w.IsSpawnTick() {
		w.SpawnMobiles()
	}
}

func (w *World) IsSpawnTick() bool {
	return w.Ticks%w.SpawnTicks == 0
}

func (w *World) IsBattleTick() bool {
	return w.Ticks%w.BattleTicks == 0
}

func (w *World) SpawnMobiles() {
	for roomUUID, rms := range w.RoomMobiles {
		room, ok := w.Rooms[roomUUID]
		if !ok {
			continue
		}
		for _, rm := range rms {
			count := room.MobileCount(rm.MobileUUID)
			for diff := rm.Count - count; diff > 0; diff-- {
				mobile := w.Mobiles.Spawn(rm.MobileUUID)
				err := room.Enter(mobile)
				if err != nil {
					log.Fatalf("Cannot spawn enter mobile %s in room %s", mobile.UUID, room.UUID)
				}
				mobile.Room = room
				mobile.Restore()
			}
		}
	}
}

type Message struct {
	FirstPerson         *Character
	SecondPerson        *Character
	FirstPersonMessage  string
	SecondPersonMessage string
	ThirdPersonMessage  string
}
