package main

import (
	"encoding/json"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/world"
	"gopkg.in/yaml.v3"
	"log"
	"os"
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

type YAMLItem struct {
	UUID          string  `yaml:"UUID"`
	Name          string  `yaml:"Name"`
	Type          string  `yaml:"Type"`
	Slot          string  `yaml:"Slot"`
	DamageType    string  `yaml:"DamageType"`
	MinimumDamage int     `yaml:"MinimumDamage"`
	MaximumDamage int     `yaml:"MaximumDamage"`
	CriticalRate  float64 `yaml:"CriticalRate"`
	CriticalBonus float64 `yaml:"CriticalBonus"`
	Modifiers     []struct {
		Type  string `yaml:"Type"`
		Value int    `yaml:"Value"`
	} `yaml:"Modifiers"`
	Keywords []string `yaml:"Keywords"`
}

type RawArea struct {
	Items []YAMLItem `yaml:"Items"`
}

func buildArea(data []byte) RawArea {
	area := RawArea{}
	err := yaml.Unmarshal(data, &area)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, item := range area.Items {
		validateItem(item)
	}
	return area
}

func validateItem(item YAMLItem) {
	if item.UUID == "" {
		log.Fatalf("Item has no UUID %v", item)
	}
	if item.Name == "" {
		log.Fatalf("Item has no Name %v", item)
	}
	if item.Type == "" {
		log.Fatalf("Item has no Type %v", item)
	}
	if item.Type != "Weapon" && item.Slot == "" {
		log.Fatalf("Item has no Slot %v", item)
	}
	for _, modifier := range item.Modifiers {
		if modifier.Type == "" {
			log.Fatalf("Item modifier has no Type %v", item)
		}
		if modifier.Value == 0 {
			log.Fatalf("Item value has no Value %v", item)
		}
	}
	if len(item.Keywords) == 0 {

		log.Fatalf("Item has no keywords %v", item)
	}
	for _, keyword := range item.Keywords {
		if keyword == "" {
			log.Fatalf("Item keyword is empty %v", item)
		}
	}
}

func buildItems(w *world.World) {
	itemFilePath := "data/item.jsonl"
	data, err := os.ReadFile(itemFilePath)
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
	Exits       struct {
		East  string `json:"east"`
		North string `json:"north"`
		South string `json:"south"`
		West  string `json:"west"`
	} `json:"exits"`
	Size    int `json:"size"`
	Mobiles []struct {
		UUID  string `json:"UUID"`
		Count int    `json:"count"`
	} `json:"mobiles"`
}

func buildRooms(w *world.World) {
	roomFilePath := "data/room.jsonl"
	data, err := os.ReadFile(roomFilePath)
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
		room.Exits.East = rr.Exits.East
		room.Exits.North = rr.Exits.North
		room.Exits.South = rr.Exits.South
		room.Exits.West = rr.Exits.West
	}
}

type RawPlayer struct {
	UUID     string `json:"UUID"`
	Name     string `json:"Name"`
	RoomUUID string `json:"RoomUUID"`
	Essence  int    `json:"Essence"`
	Power    int    `json:"Power"`
	Agility  int    `json:"Agility"`
	Insight  int    `json:"Insight"`
	Will     int    `json:"Will"`
	Health   struct {
		Current int `json:"Current"`
	} `json:"Health"`
	Spirit struct {
		Current int `json:"Current"`
	} `json:"Spirit"`
	IsAggressive bool `json:"IsAggressive"`
	IsSocial     bool `json:"IsSocial"`
	Gear         struct {
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
	Skills    struct {
		Barrier int `json:"Barrier"`
		Bash    int `json:"Bash"`
		Dagger  int `json:"Dagger"`
		Evasion int `json:"Evasion"`
		Haste   int `json:"Haste"`
		Parry   int `json:"Parry"`
		Spear   int `json:"Spear"`
		Sword   int `json:"Sword"`
	} `json:"Skills"`
}

func buildPlayers(w *world.World) {
	playerFilePath := "data/player.jsonl"
	data, err := os.ReadFile(playerFilePath)
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

		c.Essence = rp.Essence
		c.Core.Power.Base = rp.Power
		c.Core.Agility.Base = rp.Agility
		c.Core.Insight.Base = rp.Insight
		c.Core.Will.Base = rp.Will

		c.Health.Current = rp.Health.Current
		c.Spirit.Current = rp.Spirit.Current

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
		c.Skills.Bash.Base = rp.Skills.Bash
		c.Skills.Barrier.Base = rp.Skills.Barrier
		c.Skills.Dagger.Base = rp.Skills.Dagger
		c.Skills.Evasion.Base = rp.Skills.Evasion
		c.Skills.Haste.Base = rp.Skills.Haste
		c.Skills.Parry.Base = rp.Skills.Parry
		c.Skills.Spear.Base = rp.Skills.Spear
		c.Skills.Sword.Base = rp.Skills.Sword
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
	data, err := os.ReadFile(mobileFilePath)
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

		c.Essence = rp.Essence
		c.Core.Power.Base = rp.Power
		c.Core.Agility.Base = rp.Agility
		c.Core.Insight.Base = rp.Insight
		c.Core.Will.Base = rp.Will

		c.Health.Current = rp.Health.Current
		c.Spirit.Current = rp.Spirit.Current

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

		c.IsPlayer = false
		c.IsAggressive = rp.IsAggressive
		c.IsSocial = rp.IsSocial

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
