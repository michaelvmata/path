package main

import (
	"testing"
)

func TestBuild(t *testing.T) {
	rooms := build()
	expectedRooms := 2
	if len(rooms) != expectedRooms {
		t.Fatalf("Rooms count(%d) expected (%d)", len(rooms), expectedRooms)
	}
	for _, room := range rooms {
		if room.uuid == "" {
			t.Fatalf("Room uuid is empty")
		}
		if room.size == 0 {
			t.Fatalf("Room size is zero")
		}
	}
}
