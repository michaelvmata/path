package main

import (
	"github.com/michaelvmata/path/session"
	"testing"
)

func TestDetermineCommand(t *testing.T) {
	input := Noop{}.Label()
	world := build()
	s := session.New()
	c := determineCommand(input)
	go func(outgoing chan string) {
		<-outgoing
	}(s.Outgoing)
	c.Execute(world, s, input)
}
