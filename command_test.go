package main

import (
	"github.com/michaelvmata/path/world"
	"log"
	"testing"
)

func TestDetermineCommand(t *testing.T) {
	input := Noop{}.Label()
	world := build()
	c := determineCommand(input)
	ctx := Context{World: world, Player: world.Players["gaigen"], Raw: input}
	c.Execute(ctx)
}

func TestHasteCommand(t *testing.T) {
	player := world.NewPlayer("Test UUID", "Test Handle")
	ctx := Context{Player: player}
	h := Haste{}
	h.Execute(ctx)
	if len(player.Buffs) != 1 {
		t.Fatalf("Haste command didn't apply haste.")
	}
}

func TestInvest(t *testing.T) {
	previous := 0
	for i := 0; i < 100; i++ {
		cost := essenceCost(i)
		if cost <= previous {
			log.Fatalf("No incremental increase in cost %d, %d", i, cost)
		}
	}
}
