package main

import (
	"fmt"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/session"
	"github.com/michaelvmata/path/symbols"
	"github.com/michaelvmata/path/world"
	"strings"
)

type Drop struct{}

func (d Drop) Execute(w *world.World, s *session.Session, raw string) {
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 1 {
		s.Outgoing <- "Drop what?"
		return
	}
	keyword := parts[1]
	player := w.Players[s.PlayerName]
	i := player.Discard(keyword)
	if i == nil {
		s.Outgoing <- fmt.Sprintf("You don't have '%s'.", keyword)
		return
	}
	if err := player.Room.Accept(i); err != nil {
		s.Outgoing <- fmt.Sprintf("You can't drop %s.", i.Name())
		player.Receive(i)
		return
	}
	s.Outgoing <- fmt.Sprintf("You drop %s.", i.Name())
}

func (d Drop) Label() string {
	return "drop"
}

type Gear struct{}

func (g Gear) SafeName(i item.Item) string {
	if item.IsNil(i) {
		return ""
	}
	return i.Name()
}

func (g Gear) Execute(w *world.World, s *session.Session, raw string) {
	gear := w.Players[s.PlayerName].Gear
	s.Outgoing <- ""
	s.Outgoing <- fmt.Sprintf("     [Head]: %s", g.SafeName(gear.Head))
	s.Outgoing <- fmt.Sprintf("     [Neck]: %s", g.SafeName(gear.Neck))
	s.Outgoing <- fmt.Sprintf("     [Body]: %s", g.SafeName(gear.Body))
	s.Outgoing <- fmt.Sprintf("     [Arms]: %s", g.SafeName(gear.Arms))
	s.Outgoing <- fmt.Sprintf("    [Hands]: %s", g.SafeName(gear.Hands))
	s.Outgoing <- fmt.Sprintf("    [Waist]: %s", g.SafeName(gear.Waist))
	s.Outgoing <- fmt.Sprintf("     [Legs]: %s", g.SafeName(gear.Legs))
	s.Outgoing <- fmt.Sprintf("     [Feet]: %s", g.SafeName(gear.Feet))
	s.Outgoing <- fmt.Sprintf("     Wrist]: %s", g.SafeName(gear.Wrist))
	s.Outgoing <- fmt.Sprintf("  [Fingers]: %s", g.SafeName(gear.Fingers))
	s.Outgoing <- fmt.Sprintf(" [Off Hand]: %s", g.SafeName(gear.OffHand))
	s.Outgoing <- fmt.Sprintf("[Main Hand]: %s", g.SafeName(gear.MainHand))
	s.Outgoing <- ""
}

func (g Gear) Label() string {
	return "gear"
}

type Get struct{}

func (g Get) Execute(w *world.World, s *session.Session, raw string) {
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 1 {
		s.Outgoing <- "Get what?"
		return
	}
	keyword := parts[1]
	player := w.Players[s.PlayerName]
	i, err := player.Room.PickupItem(keyword)
	if err != nil {
		s.Outgoing <- fmt.Sprintf("You don't see '%s'.", keyword)
		return
	}
	if err := player.Receive(i); err != nil {
		s.Outgoing <- fmt.Sprintf("You can't get %s.", i.Name())
		player.Room.Accept(i)
		return
	} else {
		s.Outgoing <- fmt.Sprintf("You get %s.", i.Name())
	}
}

func (g Get) Label() string {
	return "get"
}

type Inventory struct{}

func (i Inventory) Execute(w *world.World, s *session.Session, raw string) {
	s.Outgoing <- ""
	player := w.Players[s.PlayerName]
	for _, i := range player.Inventory.Items {
		s.Outgoing <- i.Name()
	}
	s.Outgoing <- ""
}

func (i Inventory) Label() string {
	return "inventory"
}

type Look struct{}

func (l Look) Execute(w *world.World, s *session.Session, raw string) {
	s.Outgoing <- ""
	s.Outgoing <- w.Players[s.PlayerName].Room.Describe()
	s.Outgoing <- ""
}

func (l Look) Label() string {
	return "look"
}

type Remove struct{}

func (r Remove) Execute(w *world.World, s *session.Session, raw string) {
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 1 {
		s.Outgoing <- "Remove what?"
		return
	}
	keyword := parts[1]
	player := w.Players[s.PlayerName]
	i := player.Gear.Remove(keyword)
	if i == nil {
		s.Outgoing <- fmt.Sprintf("You don't have a '%s'", keyword)
		return
	}
	player.Inventory.AddItem(i)
	s.Outgoing <- fmt.Sprintf("You remove a %s", i.Name())
}

func (r Remove) Label() string {
	return "remove"
}

type Score struct{}

func (sc Score) Execute(w *world.World, s *session.Session, raw string) {
	p := w.Players[s.PlayerName]
	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("<red>%s<reset> Health %d(%d)+%d", symbols.HEART, p.Health.Current, p.Health.Maximum, p.Health.RecoverRate))
	parts = append(parts, fmt.Sprintf("<green>%s<reset> Spirit %d(%d)+%d", symbols.TWELVE_STAR, p.Spirit.Current, p.Spirit.Maximum, p.Spirit.RecoverRate))
	s.Outgoing <- ""
	s.Outgoing <- strings.Join(parts, "   ")
	s.Outgoing <- ""
	s.Outgoing <- p.Core.Describe()
	s.Outgoing <- ""
	s.Outgoing <- p.Skills.Describe()
	s.Outgoing <- ""
}

func (sc Score) Label() string {
	return "score"
}

type Typo struct{}

func (t Typo) Execute(w *world.World, s *session.Session, raw string) {
	s.Outgoing <- "The typo monster strikes again"
}

func (t Typo) Label() string {
	return "typo"
}

type Noop struct{}

func (n Noop) Execute(w *world.World, s *session.Session, raw string) {
	s.Outgoing <- ""
}

func (n Noop) Label() string {
	return ""
}

type Wear struct{}

func (wr Wear) Execute(w *world.World, s *session.Session, raw string) {
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 1 {
		s.Outgoing <- "Wear what?"
		return
	}
	player := w.Players[s.PlayerName]
	index := player.Inventory.IndexOfItem(parts[1])
	if index == -1 {
		s.Outgoing <- fmt.Sprintf("You don't have a '%s'", parts[1])
		return
	}
	i := player.Inventory.RemItemAtIndex(index)
	previous, err := player.Gear.Equip(i)
	if err != nil {
		s.Outgoing <- fmt.Sprintf("You can't wear %s.", i.Name())
		player.Inventory.AddItem(i)
	} else {
		s.Outgoing <- fmt.Sprintf("You wear %s", i.Name())
	}
	if !item.IsNil(previous) {
		player.Inventory.AddItem(previous)
	}
}

func (wr Wear) Label() string {
	return "wear"
}

type Executor interface {
	Execute(w *world.World, s *session.Session, raw string)
	Label() string
}

var commands = buildCommands()

func buildCommands() map[string]Executor {
	commands := []Executor{
		Drop{},
		Gear{},
		Get{},
		Inventory{},
		Look{},
		Noop{},
		Remove{},
		Score{},
		Typo{},
		Wear{},
	}
	aliases := make(map[string]Executor)
	for _, command := range commands {
		for i := range command.Label() {
			alias := command.Label()[:i+1]
			aliases[alias] = command
		}
	}
	return aliases
}

func determineCommand(raw string) Executor {
	rawCmd := strings.SplitN(raw, " ", 2)[0]
	command, ok := commands[rawCmd]
	if !ok {
		return commands[Typo{}.Label()]
	}

	return command
}
