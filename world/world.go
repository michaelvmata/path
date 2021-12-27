package world

import (
	"errors"
	"fmt"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/modifiers"
	"github.com/michaelvmata/path/session"
	"github.com/michaelvmata/path/skills"
	"github.com/michaelvmata/path/stats"
	"strings"
)

type Player struct {
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

	Attacking map[string]*Player
}

func (c Player) Clone(target Player) {
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
	c.Attacking = make(map[string]*Player)
}

func NewPlayer(UUID string, handle string) *Player {
	return &Player{
		UUID: UUID,
		Name: handle,

		keywords: strings.Fields(strings.ToLower(handle)),

		Health: stats.Line{},
		Spirit: stats.Line{},

		Core:   stats.NewCore(),
		Skills: skills.NewSkills(),

		Gear:      item.NewGear(),
		Inventory: item.NewContainer(10),

		Attacking: make(map[string]*Player),
	}
}

func (c Player) HasKeyword(target string) bool {
	target = strings.ToLower(target)
	for _, keyword := range c.keywords {
		if keyword == target {
			return true
		}
	}
	return false
}

func (c *Player) StartAttacking(defender *Player) {
	c.Attacking[defender.UUID] = defender
}

func (c *Player) StopAttacking(defender *Player) {
	delete(c.Attacking, defender.UUID)
}

func (c Player) IsAttacking(defender *Player) bool {
	_, ok := c.Attacking[defender.UUID]
	return ok
}

func (c Player) IsFighting() bool {
	return len(c.Attacking) > 0
}

func (c *Player) Discard(keyword string) item.Item {
	index := c.Inventory.IndexOfItem(keyword)
	if index == -1 {
		return nil
	}
	i := c.Inventory.RemItemAtIndex(index)
	return i
}

func (c *Player) Receive(i item.Item) error {
	if err := c.Inventory.AddItem(i); err != nil {
		return errors.New("player can't carry item")
	}
	return nil
}

func (c *Player) Move(r *Room) {
	c.Room = r
}

func (c *Player) ApplyModifiers(mods []modifiers.Modifier) {
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

func (c *Player) ApplyItemModifiers(i item.Item) {
	if !item.IsNil(i) {
		c.ApplyModifiers(i.Modifiers())
	}
}

func (c *Player) CalculateModifiers() {
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

func (c Player) IsDead() bool {
	return c.Health.Current > 0
}

func (c *Player) Update(tick int) {
	c.CalculateModifiers()
	// Adjust lines from core stats

	c.Health.Maximum = c.Core.Power.Value() * 100
	c.Health.EnforceMaximum()
	c.Spirit.Maximum = c.Core.Will.Value() * 100
	c.Spirit.EnforceMaximum()

	if tick > 0 {
		c.Health.Recover()
		c.Spirit.Recover()
	}
}

func (c Player) Show(message string, args ...interface{}) {
	if c.Session != nil {
		c.Session.Outgoing <- fmt.Sprintf(message, args...)
	}
}

func (c Player) Describe() string {
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
	Players     []*Player
	Items       item.Container
	Size        int
}

func NewRoom(uuid string, name string, description string, size int) *Room {
	room := Room{
		UUID:        uuid,
		name:        name,
		description: description,
		Size:        size,
		Players:     make([]*Player, 0, size),
		Items:       item.NewContainer(100),
	}
	return &room
}

func (r *Room) Describe(firstPerson *Player) string {
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

func (r *Room) Enter(c *Player) error {
	if r.IsFull() {
		return errors.New("room is full")
	}
	if r.IndexOfPlayer(c) != -1 {
		return errors.New("player already in room")
	}
	r.Players = append(r.Players, c)
	return nil
}

func (r *Room) Exit(c *Player) error {
	i := r.IndexOfPlayer(c)
	if i == -1 {
		return errors.New("player not in room")
	}
	copy(r.Players[i:], r.Players[:i+1])
	length := len(r.Players) - 1
	r.Players[length] = nil
	r.Players = r.Players[:length]
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

func (r *Room) IndexOfPlayer(target *Player) int {
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
	Prototypes map[string]Player
	Instances  []*Player
}

func (m *Mobiles) AddPrototype(p Player) {
	m.Prototypes[p.UUID] = p
}

func (m *Mobiles) Spawn(UUID string) *Player {
	prototype := m.Prototypes[UUID]
	mobile := NewPlayer(UUID, prototype.Name)
	mobile.Clone(prototype)
	m.Instances = append(m.Instances, mobile)
	return mobile
}

type World struct {
	Players     map[string]*Player
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
		Players: make(map[string]*Player),
		Mobiles: Mobiles{
			Prototypes: make(map[string]Player),
			Instances:  make([]*Player, 0),
		},
		Rooms:       make(map[string]*Room),
		RoomMobiles: make(map[string][]RoomMobile, 0),
		Items:       make(map[string]item.Item),
		SpawnTicks:  60,
		BattleTicks: 3,
	}
	return &w
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
			}
		}
	}
}
