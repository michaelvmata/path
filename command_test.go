package main

import "testing"

func TestDetermineCommand(t *testing.T) {
	input := Noop{}.Label()
	world := build()
	session := NewSession()
	session.player = world.Players["gaigen"]
	c := determineCommand(input)
	go func(outgoing chan string) {
		<-outgoing
	}(session.outgoing)
	c.Execute(world, session, input)
}
