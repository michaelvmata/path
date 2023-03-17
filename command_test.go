package main

import (
	"github.com/michaelvmata/path/world"
	"log"
	"testing"
)

func TestDetermineCommand(t *testing.T) {
	input := Noop{}.Label()
	world := build("data/areas")
	ctx := Context{World: world, Player: world.Players["gaigen"], Raw: input}
	c := determineCommand(input, ctx)
	c.Execute(ctx)
}

func TestBashCommand(t *testing.T) {
	world := build("data/areas")
	world.SpawnMobiles()

	// Missing target
	ctx := Context{World: world, Player: world.Players["gaigen"], Raw: "bash target"}
	b := Bash{}
	b.Execute(ctx)

	// No taret
	ctx.Raw = "bash"
	b.Execute(ctx)

	// Happy case
	ctx.Raw = "bash dummy"
	b.Execute(ctx)
}

func TestHasteCommand(t *testing.T) {
	player := world.NewPlayer("Test UUID", "Test Handle")
	ctx := Context{Player: player}
	h := Haste{}
	h.Execute(ctx)
	if len(player.Buffs) > 1 || player.Skills.Haste.IsAvailable() {
		t.Fatalf("Haste command applied without investing.")
	}
	player.Skills.Haste.Increment()
	h.Execute(ctx)
	if len(player.Buffs) < 1 || !player.Skills.Haste.IsAvailable() {
		t.Fatalf("Haste command didn't apply haste.")
	}
	h.Execute(ctx)
	if len(player.Buffs) > 1 {
		t.Fatalf("Haste command applied twice.")
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
