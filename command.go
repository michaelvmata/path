package main

type Look struct{}

func (l Look) Execute(w *World, s *Session, raw string) {
	s.outgoing <- ""
	s.outgoing <- s.player.Room.Describe()
	s.outgoing <- ""
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
	Look{}.Label(): Look{},
	Noop{}.Label(): Noop{},
	Typo{}.Label(): Typo{},
}

func determineCommand(raw string) Executor {
	command, ok := commands[raw]
	if !ok {
		return commands[Typo{}.Label()]
	}

	return command
}
