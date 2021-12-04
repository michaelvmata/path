package main

import (
	"encoding/json"
	"github.com/michaelvmata/path/items"
	"io/ioutil"
	"log"
	"strings"
)

type RawItem struct {
	UUID      string `json:"UUID"`
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	Slot      string `json:"Slot"`
	Modifiers []struct {
		Type  string `json:"Type"`
		Value int    `json:",Value"`
	} `json:"Modifiers"`
}

func buildItems(world *World) {
	itemFilePath := "data/item.ndjson"
	data, err := ioutil.ReadFile(itemFilePath)
	if err != nil {
		log.Fatalf("Error opening room %s", itemFilePath)
	}
	d := json.NewDecoder(strings.NewReader(string(data)))
	for d.More() {
		r := RawItem{}
		err := d.Decode(&r)
		if err != nil {
			log.Fatalf("Error parsing %s", data)
		}
		i := item.NewItem(r.UUID, r.Name, r.Type, r.Slot)
		for _, rm := range r.Modifiers {
			i.AddModifier(rm.Type, rm.Value)
		}
		world.Items[i.UUID] = i
	}
}

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
	Name      string `json:"name"`
	RoomUUID  string `json:"roomUuid"`
	Power     int    `json:"Power"`
	Agility   int    `json:"Agility"`
	Endurance int    `json:"Endurance"`
	Talent    int    `json:"Talent"`
	Insight   int    `json:"Insight"`
	Will      int    `json:"Will"`
	Health    struct {
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
	Gear struct {
		Head     string `json:"Head"`
		Neck     string `json:"Neck"`
		Body     string `json:"Body"`
		Arms     string `json:"Arms"`
		Hands    string `json:"Hands"`
		Waist    string `json:"Waist"`
		Legs     string `json:"Legs"`
		Feet     string `json:"Feet"`
		Wrist    string `json:"Wrist"`
		Fingers  string `json:"Fingers"`
		OffHand  string `json:"Offhand"`
		MainHand string `json:"MainHand"`
	} `json:"Gear"`
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
		c.Health.RecoverRate = rp.Health.Recover

		c.Energy.Current = rp.Energy.Current
		c.Energy.Maximum = rp.Energy.Maximum
		c.Energy.RecoverRate = rp.Energy.Recover

		c.Spirit.Current = rp.Spirit.Current
		c.Spirit.Maximum = rp.Spirit.Maximum
		c.Spirit.RecoverRate = rp.Spirit.Recover

		c.Core.Power.Base = rp.Power
		c.Core.Agility.Base = rp.Agility
		c.Core.Endurance.Base = rp.Endurance
		c.Core.Talent.Base = rp.Talent
		c.Core.Insight.Base = rp.Insight
		c.Core.Will.Base = rp.Will

		if rp.Gear.Head != "" {
			if i, ok := world.Items[rp.Gear.Head]; ok {
				c.Gear.Head = i
			}
		}
		if rp.Gear.Neck != "" {
			if i, ok := world.Items[rp.Gear.Neck]; ok {
				c.Gear.Neck = i
			}
		}
		if rp.Gear.Body != "" {
			if i, ok := world.Items[rp.Gear.Body]; ok {
				c.Gear.Body = i
			}
		}
		if rp.Gear.Arms != "" {
			if i, ok := world.Items[rp.Gear.Arms]; ok {
				c.Gear.Arms = i
			}
		}
		if rp.Gear.Hands != "" {
			if i, ok := world.Items[rp.Gear.Hands]; ok {
				c.Gear.Hands = i
			}
		}
		if rp.Gear.Waist != "" {
			if i, ok := world.Items[rp.Gear.Waist]; ok {
				c.Gear.Waist = i
			}
		}
		if rp.Gear.Legs != "" {
			if i, ok := world.Items[rp.Gear.Legs]; ok {
				c.Gear.Legs = i
			}
		}
		if rp.Gear.Feet != "" {
			if i, ok := world.Items[rp.Gear.Feet]; ok {
				c.Gear.Feet = i
			}
		}
		if rp.Gear.Wrist != "" {
			if i, ok := world.Items[rp.Gear.Wrist]; ok {
				c.Gear.Wrist = i
			}
		}
		if rp.Gear.Fingers != "" {
			if i, ok := world.Items[rp.Gear.Fingers]; ok {
				c.Gear.Fingers = i
			}
		}
		if rp.Gear.OffHand != "" {
			if i, ok := world.Items[rp.Gear.OffHand]; ok {
				c.Gear.OffHand = i
			}
		}
		if rp.Gear.MainHand != "" {
			if i, ok := world.Items[rp.Gear.MainHand]; ok {
				c.Gear.MainHand = i
			}
		}
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
	buildItems(world)
	buildRooms(world)
	buildPlayers(world)
	return world
}
