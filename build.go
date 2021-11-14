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

func build() map[string]*Room {
	roomFilePath := "data/room.ndjson"
	data, err := ioutil.ReadFile(roomFilePath)
	if err != nil {
		log.Fatalf("Error opening room %s", roomFilePath)
	}
	rooms := make(map[string]*Room)
	d := json.NewDecoder(strings.NewReader(string(data)))
	for d.More() {
		rr := RawRoom{}
		err := d.Decode(&rr)
		if err != nil {
			log.Fatalf("Error parsing %s", data)
		}
		room := NewRoom(rr.UUID, rr.Name, rr.Description, rr.Size)
		rooms[room.uuid] = room
	}
	return rooms
}
