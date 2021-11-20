package main

type Look struct{}

func (l Look) Execute(w *World, s *Session, raw string) {
	s.outgoing <- "look"
}

type Executor interface {
	Execute(w *World, s *Session, raw string)
}

var commands = map[string]Executor{
	"look": Look{},
}

func determineCommand(raw string) Executor {
	return commands[raw]
}
