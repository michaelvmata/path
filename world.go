package main

import (
	"errors"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/modifiers"
	"github.com/michaelvmata/path/skills"
	"github.com/michaelvmata/path/stats"
	"strings"
)

type Player struct {
	Name   string
	Room   *Room
	Health stats.Line
	Spirit stats.Line
	Energy stats.Line

	Core   *stats.Core
	Skills *skills.Skills

	Gear      *item.Gear
	Inventory item.Container
}

func NewPlayer(handle string) *Player {
	return &Player{
		Name:   handle,
		Health: stats.Line{},
		Spirit: stats.Line{},
		Energy: stats.Line{},

		Core:   stats.NewCore(),
		Skills: skills.NewSkills(),

		Gear:      item.NewGear(),
		Inventory: item.NewContainer(10),
	}
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
		case modifiers.Endurance:
			c.Core.Endurance.Modify(mod.Value)
		case modifiers.Talent:
			c.Core.Talent.Modify(mod.Value)
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

func (c *Player) ApplyItemModifiers(i *item.Armor) {
	if i != nil {
		c.ApplyModifiers(i.Modifiers)
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

func (c *Player) Update(tick bool) {
	c.CalculateModifiers()
	// Adjust lines from core stats
	c.Health.Maximum = c.Core.Endurance.Value() * 100
	c.Health.EnforceMaximum()
	c.Spirit.Maximum = c.Core.Insight.Value() * 100
	c.Spirit.EnforceMaximum()
	c.Energy.Maximum = (c.Core.Power.Value() + c.Core.Will.Value()) * 100
	c.Energy.EnforceMaximum()

	if tick {
		c.Health.Recover()
		c.Spirit.Recover()
		c.Energy.Recover()
	}
}

type Room struct {
	uuid        string
	name        string
	description string
	players     []*Player
	Items       item.Container
	size        int
}

func NewRoom(uuid string, name string, description string, size int) *Room {
	room := Room{
		uuid:        uuid,
		name:        name,
		description: description,
		size:        size,
		players:     make([]*Player, 0, size),
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
	return r.size == len(r.players)
}

func (r *Room) Enter(c *Player) error {
	if r.IsFull() {
		return errors.New("room is full")
	}
	if r.IndexOfPlayer(c) != -1 {
		return errors.New("player already in room")
	}
	r.players = append(r.players, c)
	return nil
}

func (r *Room) Exit(c *Player) error {
	i := r.IndexOfPlayer(c)
	if i == -1 {
		return errors.New("player not in room")
	}
	copy(r.players[i:], r.players[:i+1])
	length := len(r.players) - 1
	r.players[length] = nil
	r.players = r.players[:length]
	return nil
}

func (r *Room) IndexOfPlayer(target *Player) int {
	for i, p := range r.players {
		if p == target {
			return i
		}
	}
	return -1
}

type World struct {
	Players map[string]*Player
	Rooms   map[string]*Room
	Items   map[string]*item.Armor
}

func NewWorld() *World {
	w := World{
		Players: make(map[string]*Player),
		Rooms:   make(map[string]*Room),
		Items:   make(map[string]*item.Armor),
	}
	return &w
}
