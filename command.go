package main

import (
	"fmt"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/symbols"
	"strings"
)

type Drop struct{}

func (d Drop) Execute(w *World, s *Session, raw string) {
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 1 {
		s.outgoing <- "Drop what?"
		return
	}
	keyword := parts[1]
	i := s.player.Discard(keyword)
	if i == nil {
		s.outgoing <- fmt.Sprintf("You don't have '%s'.", keyword)
		return
	}
	if err := s.player.Room.Accept(i); err != nil {
		s.outgoing <- fmt.Sprintf("You can't drop %s.", i.Name)
		s.player.Receive(i)
		return
	}
	s.outgoing <- fmt.Sprintf("You drop %s.", i.Name)
}

func (d Drop) Label() string {
	return "drop"
}

type Gear struct{}

func (g Gear) SafeName(i *item.Item) string {
	if i == nil {
		return ""
	}
	return i.Name
}

func (g Gear) Execute(w *World, s *Session, raw string) {
	gear := s.player.Gear
	s.outgoing <- ""
	s.outgoing <- fmt.Sprintf("     [Head]: %s", g.SafeName(gear.Head))
	s.outgoing <- fmt.Sprintf("     [Neck]: %s", g.SafeName(gear.Neck))
	s.outgoing <- fmt.Sprintf("     [Body]: %s", g.SafeName(gear.Body))
	s.outgoing <- fmt.Sprintf("     [Arms]: %s", g.SafeName(gear.Arms))
	s.outgoing <- fmt.Sprintf("    [Hands]: %s", g.SafeName(gear.Hands))
	s.outgoing <- fmt.Sprintf("    [Waist]: %s", g.SafeName(gear.Waist))
	s.outgoing <- fmt.Sprintf("     [Legs]: %s", g.SafeName(gear.Legs))
	s.outgoing <- fmt.Sprintf("     [Feet]: %s", g.SafeName(gear.Feet))
	s.outgoing <- fmt.Sprintf("     Wrist]: %s", g.SafeName(gear.Wrist))
	s.outgoing <- fmt.Sprintf("  [Fingers]: %s", g.SafeName(gear.Fingers))
	s.outgoing <- fmt.Sprintf(" [Off Hand]: %s", g.SafeName(gear.OffHand))
	s.outgoing <- fmt.Sprintf("[Main Hand]: %s", g.SafeName(gear.MainHand))
	s.outgoing <- ""
}

func (g Gear) Label() string {
	return "gear"
}

type Get struct{}

func (g Get) Execute(w *World, s *Session, raw string) {
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 1 {
		s.outgoing <- "Get what?"
		return
	}
	keyword := parts[1]
	i, err := s.player.Room.PickupItem(keyword)
	if err != nil {
		s.outgoing <- fmt.Sprintf("You don't see '%s'.", keyword)
		return
	}
	if err := s.player.Receive(i); err != nil {
		s.outgoing <- fmt.Sprintf("You can't get %s.", i.Name)
		s.player.Room.Accept(i)
		return
	} else {
		s.outgoing <- fmt.Sprintf("You get %s.", i.Name)
	}
}

func (g Get) Label() string {
	return "get"
}

type Inventory struct{}

func (i Inventory) Execute(w *World, s *Session, raw string) {
	s.outgoing <- ""
	for _, i := range s.player.Inventory.Items {
		s.outgoing <- i.Name
	}
	s.outgoing <- ""
}

func (i Inventory) Label() string {
	return "inventory"
}

type Look struct{}

func (l Look) Execute(w *World, s *Session, raw string) {
	s.outgoing <- ""
	s.outgoing <- s.player.Room.Describe()
	s.outgoing <- ""
}

func (l Look) Label() string {
	return "look"
}

type Remove struct{}

func (r Remove) Execute(w *World, s *Session, raw string) {
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 1 {
		s.outgoing <- "Remove what?"
		return
	}
	keyword := parts[1]
	i := s.player.Gear.Remove(keyword)
	if i == nil {
		s.outgoing <- fmt.Sprintf("You don't have a '%s'", keyword)
		return
	}
	s.player.Inventory.AddItem(i)
	s.outgoing <- fmt.Sprintf("You remove a %s", i.Name)
}

func (r Remove) Label() string {
	return "remove"
}

type Score struct{}

func (sc Score) Execute(w *World, s *Session, raw string) {
	p := s.player
	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("<red>%s<reset> Health %d(%d)+%d", symbols.HEART, p.Health.Current, p.Health.Maximum, p.Health.RecoverRate))
	parts = append(parts, fmt.Sprintf("<green>%s<reset> Spirit %d(%d)+%d", symbols.TWELVE_STAR, p.Spirit.Current, p.Spirit.Maximum, p.Spirit.RecoverRate))
	parts = append(parts, fmt.Sprintf("<blue>%s<reset> Energy %d(%d)+%d", symbols.CIRCLED_BULLET, p.Energy.Current, p.Energy.Maximum, p.Energy.RecoverRate))
	s.outgoing <- ""
	s.outgoing <- strings.Join(parts, "   ")
	s.outgoing <- ""
	s.outgoing <- p.Core.Describe()
	s.outgoing <- ""
	s.outgoing <- p.Skills.Describe()
	s.outgoing <- ""
}

func (sc Score) Label() string {
	return "score"
}

type Typo struct{}

func (t Typo) Execute(w *World, s *Session, raw string) {
	s.outgoing <- "The typo monster strikes again"
}

func (t Typo) Label() string {
	return "typo"
}

type Noop struct{}

func (n Noop) Execute(w *World, s *Session, raw string) {
	s.outgoing <- ""
}

func (n Noop) Label() string {
	return ""
}

type Wear struct{}

func (wr Wear) Execute(w *World, s *Session, raw string) {
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) == 1 {
		s.outgoing <- "Wear what?"
		return
	}
	index := s.player.Inventory.IndexOfItem(parts[1])
	if index == -1 {
		s.outgoing <- fmt.Sprintf("You don't have a '%s'", parts[1])
		return
	}
	i := s.player.Inventory.RemItemAtIndex(index)
	previous, err := s.player.Gear.Equip(i)
	if err != nil {
		s.outgoing <- fmt.Sprintf("You can't wear %s.", i.Name)
		s.player.Inventory.AddItem(i)
	} else {
		s.outgoing <- fmt.Sprintf("You wear %s", i.Name)
	}
	if previous != nil {
		s.player.Inventory.AddItem(previous)
	}
}

func (wr Wear) Label() string {
	return "wear"
}

type Executor interface {
	Execute(w *World, s *Session, raw string)
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
