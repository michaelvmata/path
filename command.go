package main

import (
	"fmt"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/symbols"
	"strings"
)

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

type Look struct{}

func (l Look) Execute(w *World, s *Session, raw string) {
	s.outgoing <- ""
	s.outgoing <- s.player.Room.Describe()
	s.outgoing <- ""
}

func (l Look) Label() string {
	return "look"
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

type Executor interface {
	Execute(w *World, s *Session, raw string)
}

var commands = map[string]Executor{
	Gear{}.Label():  Gear{},
	Look{}.Label():  Look{},
	Noop{}.Label():  Noop{},
	Score{}.Label(): Score{},
	Typo{}.Label():  Typo{},
}

func determineCommand(raw string) Executor {
	command, ok := commands[raw]
	if !ok {
		return commands[Typo{}.Label()]
	}

	return command
}
