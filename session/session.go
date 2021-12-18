package session

type Session struct {
	PlayerName string
	Incoming   chan string
	Outgoing   chan string
}

func New() *Session {
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

func (s *Session) Send(messages []string) {
	for _, message := range messages {
		s.Outgoing <- message
	}
}
