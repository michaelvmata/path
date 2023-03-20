package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	world := build("data/areas")
	expectedRooms := 2
	if len(world.Rooms) < expectedRooms {
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

func TestSavePlayers(t *testing.T) {
	world := build("data/areas")
	savePlayers(world.Players)
}

func TestArea(t *testing.T) {
	absPath, err := filepath.Abs("data/areas/default.yaml")
	if err != nil {
		t.Fatalf("Test path error")
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Error reading YAML items file")
	}
	area := buildArea(data)
	if len(area.Items) == 0 {
		t.Fatalf("Failed to parse items")
	}
}
