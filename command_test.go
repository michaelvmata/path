package main

import "testing"

func TestDetermineCommand(t *testing.T) {
	input := Look{}.Label()
	world := build()
	session := NewSession()
	session.player = world.Players["gaigen"]
	c := determineCommand(input)
	go func(outgoing chan string) {
		output := <-outgoing
		if output == "" {
			t.Fatalf("Got empty output")
		}
	}(session.outgoing)
	c.Execute(world, session, input)
}
