package session

import (
	"log"
	"testing"
)

func TestSession(t *testing.T) {
	s := NewSession()
	if r := s.Receive(); r != "" {
		log.Fatalf("Unexpected result %s", r)
	}
}
