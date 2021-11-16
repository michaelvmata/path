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
		if room.uuid == "" {
			t.Fatalf("Room uuid is empty")
		}
		if room.size == 0 {
			t.Fatalf("Room size is zero")
		}
	}
}
