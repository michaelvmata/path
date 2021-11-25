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

type RawPlayer struct {
	RoomUUID string `json:"RoomUUID"`
	Name     string `json:"Name"`
	Health   struct {
		Maximum int `json:"Maximum"`
		Current int `json:"Current"`
		Recover int `json:"Recover"`
	} `json:"Health"`
	Spirit struct {
		Maximum int `json:"Maximum"`
		Current int `json:"Current"`
		Recover int `json:"Recover"`
	} `json:"Spirit"`
	Energy struct {
		Maximum int `json:"Maximum"`
		Current int `json:"Current"`
		Recover int `json:"Recover"`
	} `json:"Energy"`
}

func buildPlayers(world *World) {
	playerFilePath := "data/player.ndjson"
	data, err := ioutil.ReadFile(playerFilePath)
	if err != nil {
		log.Fatalf("Error opening players %s", playerFilePath)
	}
	d := json.NewDecoder(strings.NewReader(string(data)))
	for d.More() {
		rp := RawPlayer{}
		err := d.Decode(&rp)
		if err != nil {
			log.Fatalf("Error parsing %s", data)
		}
		c := NewPlayer(rp.Name)
		c.Health.Current = rp.Health.Current
		c.Health.Maximum = rp.Health.Maximum
		c.Health.Recover = rp.Health.Recover

		c.Energy.Current = rp.Energy.Current
		c.Energy.Maximum = rp.Energy.Maximum
		c.Energy.Recover = rp.Energy.Recover

		c.Spirit.Current = rp.Spirit.Current
		c.Spirit.Maximum = rp.Spirit.Maximum
		c.Spirit.Recover = rp.Spirit.Recover

		world.Players[c.Name] = c
		if room, ok := world.Rooms[rp.RoomUUID]; ok {
			if err := room.Enter(c); err == nil {
				c.Room = room
			}
		}
	}
}

func build() *World {
	world := NewWorld()
	buildRooms(world)
	buildPlayers(world)
	return world
}
