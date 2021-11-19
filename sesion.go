package main

import "errors"

type Session struct {
	player   *Player
	incoming chan string
	outgoing chan string
}

func NewSession() *Session {
	s := Session{
		incoming: make(chan string),
		outgoing: make(chan string),
	}
	return &s
}

func (s *Session) Receive() string {
	select {
	case input := <-s.incoming:
		return input
	default:
		return ""
	}
}

func (s *Session) Process() error {
	m := s.Receive()
	if m == "" {
		return nil
	}
	if m == "quit" {
		return errors.New("time to quit")
	}
	messages := make([]string, 0)
	messages = append(messages, m)
	s.Send(messages)
	return nil
}

func (s *Session) Send(messages []string) {
	for _, message := range messages {
		s.outgoing <- message
	}
}
