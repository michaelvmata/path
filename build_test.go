package main

import (
	"testing"
)

func TestBuild(t *testing.T) {
	world := build()
	expectedRooms := 2
	if len(world.Rooms) != expectedRooms {
		t.Fatalf("Rooms count(%d) expected (%d)", len(world.Rooms), expectedRooms)
	}
	for _, room := range world.Rooms {
		if room.UUID == "" {
			t.Fatalf("Room uuid is empty")
		}
		if room.Size == 0 {
			t.Fatalf("Room size is zero")
		}
	}

	for _, player := range world.Players {
		if player.Room == nil {
			t.Fatalf("Player room not loaded")
		}
		if player.Name == "" {
			t.Fatalf("Player name empty")
		}
	}
}
