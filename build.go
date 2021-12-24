package main

import (
	"encoding/json"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/world"
	"io/ioutil"
	"log"
	"strings"
)

type RawItem struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	Type          string  `json:"Type"`
	WeaponType    string  `json:"WeaponType"`
	MinimumDamage int     `json:"MinimumDamage"`
	MaximumDamage int     `json:"MaximumDamage"`
	CriticalRate  float64 `json:"CriticalRate"`
	CriticalBonus float64 `json:"CriticalBonus"`
	Slot          string  `json:"Slot"`
	Modifiers     []struct {
		Type  string `json:"Type"`
		Value int    `json:"Value"`
	} `json:"Modifiers"`
	Keywords []string `json:"Keywords"`
}

func buildItems(w *world.World) {
	itemFilePath := "data/item.jsonl"
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
		var i item.Item
		if r.Type == "Armor" {
			i = item.NewArmor(r.UUID, r.Name, r.Slot, r.Keywords)
		} else if r.Type == "Weapon" {
			w := item.NewWeapon(r.UUID, r.Name, r.Keywords, r.WeaponType)
			if r.MaximumDamage <= r.MinimumDamage || r.MinimumDamage <= 0 {
				log.Fatalln("Invalid Maximum and Minimum Damage", r)
			}
			w.MinimumDamage = r.MinimumDamage
			w.MaximumDamage = r.MaximumDamage
			w.CriticalBonus = r.CriticalBonus
			w.CriticalRate = r.CriticalRate
			i = w
		}
		for _, rm := range r.Modifiers {
			i.AddModifier(rm.Type, rm.Value)
		}
		w.Items[i.UUID()] = i
	}
}

type RawRoom struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	Mobiles     []struct {
		UUID  string `json:"UUID"`
		Count int    `json:"count"`
	} `json:"mobiles"`
}

func buildRooms(w *world.World) {
	roomFilePath := "data/room.jsonl"
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
		room := world.NewRoom(rr.UUID, rr.Name, rr.Description, rr.Size)
		w.Rooms[room.UUID] = room
		for _, mc := range rr.Mobiles {
			w.RoomMobiles[room.UUID] = append(w.RoomMobiles[room.UUID], world.NewRoomMobile(mc.UUID, mc.Count))
		}
	}
}

type RawPlayer struct {
	UUID     string `json:"UUID"`
	Name     string `json:"name"`
	RoomUUID string `json:"roomUuid"`
	Power    int    `json:"Power"`
	Agility  int    `json:"Agility"`
	Insight  int    `json:"Insight"`
	Will     int    `json:"Will"`
	Health   struct {
		Current int `json:"Current"`
		Recover int `json:"Recover"`
	} `json:"Health"`
	Spirit struct {
		Current int `json:"Current"`
		Recover int `json:"Recover"`
	} `json:"Spirit"`
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
		OffHand  string `json:"OffHand"`
		MainHand string `json:"MainHand"`
	} `json:"Gear"`
	Inventory []string `json:"Inventory"`
}

func buildPlayers(w *world.World) {
	playerFilePath := "data/player.jsonl"
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
		c := world.NewPlayer(rp.UUID, rp.Name)

		c.Core.Power.Base = rp.Power
		c.Core.Agility.Base = rp.Agility
		c.Core.Insight.Base = rp.Insight
		c.Core.Will.Base = rp.Will

		c.Health.Current = rp.Health.Current
		c.Health.RecoverRate = rp.Health.Recover

		c.Spirit.Current = rp.Spirit.Current
		c.Spirit.RecoverRate = rp.Spirit.Recover

		if rp.Gear.Head != "" {
			if i, ok := w.Items[rp.Gear.Head]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Neck != "" {
			if i, ok := w.Items[rp.Gear.Neck]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Body != "" {
			if i, ok := w.Items[rp.Gear.Body]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Arms != "" {
			if i, ok := w.Items[rp.Gear.Arms]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Hands != "" {
			if i, ok := w.Items[rp.Gear.Hands]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Waist != "" {
			if i, ok := w.Items[rp.Gear.Waist]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Legs != "" {
			if i, ok := w.Items[rp.Gear.Legs]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Feet != "" {
			if i, ok := w.Items[rp.Gear.Feet]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Wrist != "" {
			if i, ok := w.Items[rp.Gear.Wrist]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Fingers != "" {
			if i, ok := w.Items[rp.Gear.Fingers]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.OffHand != "" {
			if i, ok := w.Items[rp.Gear.OffHand]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.MainHand != "" {
			if i, ok := w.Items[rp.Gear.MainHand]; ok {
				c.Gear.Equip(i)
			}
		}
		for _, itemUUID := range rp.Inventory {
			if i, ok := w.Items[itemUUID]; ok {
				c.Inventory.AddItem(i)
			}
		}
		c.Update(0)
		w.Players[c.Name] = c
		if room, ok := w.Rooms[rp.RoomUUID]; ok {
			if err := room.Enter(c); err == nil {
				c.Room = room
			}
		}
	}
}

func buildMobiles(w *world.World) {
	mobileFilePath := "data/mobile.jsonl"
	data, err := ioutil.ReadFile(mobileFilePath)
	if err != nil {
		log.Fatalf("Error opening mobile %s", mobileFilePath)
	}
	d := json.NewDecoder(strings.NewReader(string(data)))
	for d.More() {
		rp := RawPlayer{}
		err := d.Decode(&rp)
		if err != nil {
			log.Fatalf("Error parsing %s", data)
		}
		c := world.NewPlayer(rp.UUID, rp.Name)

		c.Core.Power.Base = rp.Power
		c.Core.Agility.Base = rp.Agility
		c.Core.Insight.Base = rp.Insight
		c.Core.Will.Base = rp.Will

		c.Health.Current = rp.Health.Current
		c.Health.RecoverRate = rp.Health.Recover

		c.Spirit.Current = rp.Spirit.Current
		c.Spirit.RecoverRate = rp.Spirit.Recover

		if rp.Gear.Head != "" {
			if i, ok := w.Items[rp.Gear.Head]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Neck != "" {
			if i, ok := w.Items[rp.Gear.Neck]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Body != "" {
			if i, ok := w.Items[rp.Gear.Body]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Arms != "" {
			if i, ok := w.Items[rp.Gear.Arms]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Hands != "" {
			if i, ok := w.Items[rp.Gear.Hands]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Waist != "" {
			if i, ok := w.Items[rp.Gear.Waist]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Legs != "" {
			if i, ok := w.Items[rp.Gear.Legs]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Feet != "" {
			if i, ok := w.Items[rp.Gear.Feet]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Wrist != "" {
			if i, ok := w.Items[rp.Gear.Wrist]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.Fingers != "" {
			if i, ok := w.Items[rp.Gear.Fingers]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.OffHand != "" {
			if i, ok := w.Items[rp.Gear.OffHand]; ok {
				c.Gear.Equip(i)
			}
		}
		if rp.Gear.MainHand != "" {
			if i, ok := w.Items[rp.Gear.MainHand]; ok {
				c.Gear.Equip(i)
			}
		}
		for _, itemUUID := range rp.Inventory {
			if i, ok := w.Items[itemUUID]; ok {
				c.Inventory.AddItem(i)
			}
		}
		c.Update(0)
		w.Mobiles.AddPrototype(*c)
	}
}

func build() *world.World {
	world := world.NewWorld()
	buildItems(world)
	buildMobiles(world)
	buildRooms(world)
	buildPlayers(world)
	return world
}
