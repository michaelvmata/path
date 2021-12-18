package main

import (
	"testing"
)

func TestDetermineCommand(t *testing.T) {
	input := Noop{}.Label()
	world := build()
	c := determineCommand(input)
	ctx := Context{World: world, Player: world.Players["gaigen"], Raw: input}
	c.Execute(ctx)
}
