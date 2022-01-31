package world

import (
	"errors"
	"fmt"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/modifiers"
	"github.com/michaelvmata/path/session"
	"github.com/michaelvmata/path/skills"
	"github.com/michaelvmata/path/stats"
	"github.com/michaelvmata/path/symbols"
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

	Attacking map[string]*Character

	Buffs   []Buff
	Stunned int
}

func (c *Character) CreditEssence(amount int) {
	c.Essence += amount
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

func (c *Character) UnapplyExpiredBuffs() {
	buffs := make([]Buff, 0)
	messages := make([]string, 0)
	for _, b := range c.Buffs {
		if !b.IsExpired() {
			buffs = append(buffs, b)
		} else {
			messages = append(messages, b.UnapplyMessage())
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

func (c *Character) Stun(length int) {
	c.Stunned += length
}

func (c *Character) ReduceStun() {
	if c.Stunned > 0 {
		c.Stunned--
	}
}

func (c Character) IsStunned() bool {
	return c.Stunned > 0
}

func (c Character) Weapon() *item.Weapon {
	if c.Gear.MainHand != nil {
		return c.Gear.MainHand
	}
	hand := item.NewWeapon("Barehand", "Barehand", []string{}, item.Crush)
	hand.MinimumDamage = 5
	hand.MaximumDamage = 10
	hand.CriticalBonus = 1
	hand.CriticalRate = 0
	return hand
}

func (c *Character) Clone(target Character) {
	c.UUID = target.UUID
	c.Name = target.Name
	c.Room = nil
	c.Session = nil

	c.Core.Power.Base = target.Core.Power.Base
	c.Core.Agility.Base = target.Core.Agility.Base
	c.Core.Insight.Base = target.Core.Insight.Base
	c.Core.Will.Base = target.Core.Will.Base

	c.Gear = item.NewGear()
	c.Inventory = item.NewContainer(10)
	c.Attacking = make(map[string]*Character)
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

		Attacking: make(map[string]*Character),
	}
}

func (c Character) ShowPrompt() {
	border := "<grey_62>>> "
	c.Show("%s%s <red>%d%s <green>%d%s %s",
		border,
		c.Name,
		c.Health.Current, symbols.HEART,
		c.Spirit.Current, symbols.TWELVE_STAR,
		border)
}

func (c Character) HasKeyword(target string) bool {
	target = strings.ToLower(target)
	for _, keyword := range c.keywords {
		if keyword == target {
			return true
		}
	}
	return false
}

func (c *Character) StartAttacking(defender *Character) {
	c.Attacking[defender.UUID] = defender
}

func (c *Character) StopAttacking(defender *Character) {
	delete(c.Attacking, defender.UUID)
}

func (c Character) IsAttacking(defender *Character) bool {
	_, ok := c.Attacking[defender.UUID]
	return ok
}

func (c Character) IsFighting() bool {
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
		if tick%5 == 0 {
			c.Health.Recover()
			c.Spirit.Recover()
		}
		for _, buff := range c.Buffs {
			buff.Update(tick)
		}
		c.UnapplyExpiredBuffs()
		c.ReduceStun()
	}
}

func (c Character) ShowNewline() {
	if c.Session != nil {
		c.Session.Outgoing <- "\n"
	}
}

func (c Character) Showln(message string, args ...interface{}) {
	if c.Session != nil {
		c.Session.Outgoing <- fmt.Sprintf(message, args...)
		c.Session.Outgoing <- "\n"
	}
}

func (c Character) Show(message string, args ...interface{}) {
	if c.Session != nil {
		c.Session.Outgoing <- fmt.Sprintf(message, args...)
	}
}

func (c Character) Describe() string {
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

type Room struct {
	UUID        string
	name        string
	description string
	Players     []*Character
	Items       item.Container
	Size        int
}

func NewRoom(uuid string, name string, description string, size int) *Room {
	room := Room{
		UUID:        uuid,
		name:        name,
		description: description,
		Size:        size,
		Players:     make([]*Character, 0, size),
		Items:       item.NewContainer(100),
	}
	return &room
}

func (r *Room) Describe(firstPerson *Character) string {
	parts := make([]string, 0)
	parts = append(parts, r.name)
	parts = append(parts, "")
	parts = append(parts, r.description)
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

func (r Room) IndexOfPlayerHandle(handle string) int {
	for i, p := range r.Players {
		if p.HasKeyword(handle) {
			return i
		}
	}
	return -1
}

func (r *Room) IndexOfPlayer(target *Character) int {
	for i, p := range r.Players {
		if p == target {
			return i
		}
	}
	return -1
}

func (r Room) MobileCount(mobileUUID string) int {
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

func (m Mobiles) IsInstance(c *Character) bool {
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
		SpawnTicks:  1,
		BattleTicks: 3,
	}
	return &w
}

func (w World) IsMobile(c *Character) bool {
	return w.Mobiles.IsInstance(c)
}

func (w *World) Update() {
	w.Ticks++
	for _, player := range w.Players {
		player.Update(w.Ticks)
	}
	if w.IsSpawnTick() {
		w.SpawnMobiles()
	}
}

func (w World) IsSpawnTick() bool {
	return w.Ticks%w.SpawnTicks == 0
}

func (w World) IsBattleTick() bool {
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
				room.Enter(mobile)
				mobile.Room = room
				mobile.Restore()
			}
		}
	}
}
