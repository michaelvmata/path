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

	Health stats.Line
	Spirit stats.Line

	Core   stats.Core
	Skills skills.Skills

	Gear      *item.Gear
	Inventory item.Container

	Attacking map[string]bool
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
	c.Attacking = make(map[string]bool)
}

func NewPlayer(UUID string, handle string) *Player {
	return &Player{
		UUID: UUID,
		Name: handle,

		Health: stats.Line{},
		Spirit: stats.Line{},

		Core:   stats.NewCore(),
		Skills: skills.NewSkills(),

		Gear:      item.NewGear(),
		Inventory: item.NewContainer(10),

		Attacking: make(map[string]bool),
	}
}

func (c *Player) StartAttacking(defender string) {
	c.Attacking[defender] = true
}

func (c *Player) StopAttacking(defender string) {
	delete(c.Attacking, defender)
}

func (c Player) IsAttacking(defender string) bool {
	_, ok := c.Attacking[defender]
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

func (c *Player) Update(tick bool) {
	c.CalculateModifiers()
	// Adjust lines from core stats

	c.Health.Maximum = c.Core.Power.Value() * 100
	c.Health.EnforceMaximum()
	c.Spirit.Maximum = c.Core.Will.Value() * 100
	c.Spirit.EnforceMaximum()

	if tick {
		c.Health.Recover()
		c.Spirit.Recover()
	}
}

func (c Player) Show(message string, args ...interface{}) {
	if c.Session != nil {
		c.Session.Outgoing <- fmt.Sprintf(message, args...)
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

func (r *Room) Describe() string {
	parts := make([]string, 0)
	parts = append(parts, r.name)
	parts = append(parts, "")
	parts = append(parts, r.description)
	parts = append(parts, "")
	for _, i := range r.Items.Items {
		parts = append(parts, i.Name())
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

func (r *Room) IndexOfPlayer(target *Player) int {
	for i, p := range r.Players {
		if p == target {
			return i
		}
	}
	return -1
}

type Mobiles struct {
	Prototypes map[string]Player
	Instances  []*Player
	count      int
}

func (m *Mobiles) AddPrototype(p Player) {
	m.Prototypes[p.UUID] = p
}

func (m *Mobiles) Spawn(UUID string) *Player {
	if len(m.Instances) == m.count {
		m.count += 1
	}
	prototype := m.Prototypes[UUID]
	mobile := NewPlayer(UUID, prototype.Name)
	mobile.Clone(prototype)
	m.Instances[m.count-1] = mobile
	return mobile
}

type World struct {
	Players map[string]*Player
	Mobiles Mobiles
	Rooms   map[string]*Room
	Items   map[string]item.Item
}

func NewWorld() *World {
	w := World{
		Players: make(map[string]*Player),
		Mobiles: Mobiles{
			Prototypes: make(map[string]Player),
			Instances:  make([]*Player, 0),
		},
		Rooms: make(map[string]*Room),
		Items: make(map[string]item.Item),
	}
	return &w
}

func (w *World) Update(tick bool) {
	for _, player := range w.Players {
		player.Update(tick)
	}
}
