package session

import (
	"errors"
)

type Session struct {
	PlayerName string
	Incoming   chan string
	Outgoing   chan string
}

func NewSession() *Session {
	s := Session{
		Incoming: make(chan string),
		Outgoing: make(chan string),
	}
	return &s
}

func (s *Session) HasPlayer() bool {
	return s.PlayerName != ""
}

func (s *Session) Receive() string {
	select {
	case input := <-s.Incoming:
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
		s.Outgoing <- message
	}
}
