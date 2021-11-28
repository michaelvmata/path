package main

import (
	"fmt"
	"github.com/michaelvmata/path/symbols"
	"strings"
)

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
	parts = append(parts, fmt.Sprintf("<red>%s<reset> Health %d(%d)+%d", symbols.HEART, p.Health.Current, p.Health.Maximum, p.Health.Recover))
	parts = append(parts, fmt.Sprintf("<green>%s<reset> Spirit %d(%d)+%d", symbols.TWELVE_STAR, p.Spirit.Current, p.Spirit.Maximum, p.Spirit.Recover))
	parts = append(parts, fmt.Sprintf("<blue>%s<reset> Energy %d(%d)+%d", symbols.CIRCLED_BULLET, p.Energy.Current, p.Energy.Maximum, p.Energy.Recover))
	s.outgoing <- ""
	s.outgoing <- strings.Join(parts, "   ")
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
