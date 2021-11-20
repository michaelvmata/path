package main

import (
	"log"
)

type Look struct{}

func (l Look) Execute(w *World, s *Session, raw string) {
	s.outgoing <- "look"
}

func (l Look) Label() string {
	return "look"

}

type Typo struct{}

func (t Typo) Execute(w *World, s *Session, raw string) {
	s.outgoing <- "The typo monster strikes again"
}

func (t Typo) Label() string {
	return "typo"
}

type Executor interface {
	Execute(w *World, s *Session, raw string)
}

var commands = map[string]Executor{
	Look{}.Label(): Look{},
	Typo{}.Label(): Typo{},
}

func determineCommand(raw string) Executor {
	command, ok := commands[raw]
	log.Println("checking command")
	if !ok {
		return commands[Typo{}.Label()]
	}

	return command
}
