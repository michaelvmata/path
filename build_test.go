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
	expectedName := "The Void"
	if rooms[0].name != "The Void" {
		t.Fatalf("Room name(%s) expected(%s)", rooms[0].name, expectedName)
	}
	expectedSize := 2
	if rooms[0].size != expectedSize {
		t.Fatalf("Room size(%d) expected(%d)", rooms[0].size, expectedSize)
	}
}
