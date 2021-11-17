package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type RawRoom struct {
	UUID        string
	Name        string
	Description string
	Size        int
}

func buildRooms(world *World) {
	roomFilePath := "data/room.ndjson"
	data, err := ioutil.ReadFile(roomFilePath)
	if err != nil {
		log.Fatalf("Error opening room %s", roomFilePath)
	}
	d := json.NewDecoder(strings.NewReader(string(data)))
	for d.More() {
		rr := RawRoom{}
		err := d.Decode(&rr)
		if err != nil {
			log.Fatalf("Error parsing %s", data)
		}
		room := NewRoom(rr.UUID, rr.Name, rr.Description, rr.Size)
		world.Rooms[room.uuid] = room
	}
}

type RawCharacter struct {
	RoomUUID string
	Name     string
}

func buildCharacters(world *World) {
	characterFilePath := "data/character.ndjson"
	data, err := ioutil.ReadFile(characterFilePath)
	if err != nil {
		log.Fatalf("Error opening character %s", characterFilePath)
	}
	d := json.NewDecoder(strings.NewReader(string(data)))
	for d.More() {
		rc := RawCharacter{}
		err := d.Decode(&rc)
		if err != nil {
			log.Fatalf("Error parsing %s", data)
		}
		c := NewCharacter(rc.Name)
		world.Characters[c.Name] = c
		if room, ok := world.Rooms[rc.RoomUUID]; ok {
			if err := room.EnterCharacter(c); err == nil {
				c.Room = room
			}
		}
	}
}

func build() *World {
	world := NewWorld()
	buildRooms(world)
	buildCharacters(world)
	return world
}
