package main

import (
	session2 "github.com/michaelvmata/path/session"
	"testing"
)

func TestDetermineCommand(t *testing.T) {
	input := Noop{}.Label()
	world := build()
	session := session2.NewSession()
	c := determineCommand(input)
	go func(outgoing chan string) {
		<-outgoing
	}(session.Outgoing)
	c.Execute(world, session, input)
}
