package main

import "testing"

func TestDetermineCommand(t *testing.T) {
	input := "look"
	world := NewWorld()
	session := NewSession()
	c := determineCommand(input)
	go func(outgoing chan string) {
		output := <-outgoing
		if output != input {
			t.Fatalf("Got output(%s) expected(%s)", output, input)
		}
	}(session.outgoing)
	c.Execute(world, session, input)
}
