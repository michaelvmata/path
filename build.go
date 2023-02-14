package main

import (
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/world"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

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

type YAMLMobile struct {
	UUID         string `yaml:"UUID"`
	Name         string `yaml:"Name"`
	Essence      int    `yaml:"Essence"`
	Power        int    `yaml:"Power"`
	Agility      int    `yaml:"Agility"`
	Will         int    `yaml:"Will"`
	Insight      int    `yaml:"Insight"`
	IsAggressive bool   `yaml:"IsAggressive"`
	IsSocial     bool   `yaml:"IsSocial"`
	Health       int    `yaml:"Health"`
	Spirit       int    `yaml:"Spirit"`
	RoomUUID     string `yaml:"RoomUUID"`
	Gear         struct {
		Head     string `yaml:"Head"`
		Neck     string `yaml:"Neck"`
		Body     string `yaml:"Body"`
		Arms     string `yaml:"Arms"`
		Hands    string `yaml:"Hands"`
		Waist    string `yaml:"Waist"`
		Legs     string `yaml:"Legs"`
		Feet     string `yaml:"Feet"`
		Wrist    string `yaml:"Wrist"`
		Fingers  string `yaml:"Fingers"`
		OffHand  string `yaml:"OffHand"`
		MainHand string `yaml:"MainHand"`
	} `yaml:"Gear"`
	Inventory []string `yaml:"Inventory"`
	Skills    struct {
		Barrier int `yaml:"Barrier"`
		Bash    int `yaml:"Bash"`
		Dagger  int `yaml:"Dagger"`
		Evasion int `yaml:"Evasion"`
		Haste   int `yaml:"Haste"`
		Parry   int `yaml:"Parry"`
		Spear   int `yaml:"Spear"`
		Sword   int `yaml:"Sword"`
	} `yaml:"Skills"`
}

type YAMLRoom struct {
	UUID        string `yaml:"UUID"`
	Name        string `yaml:"Name"`
	Description string `yaml:"Description"`
	Size        int    `yaml:"Size"`
	Exits       struct {
		East  string `yaml:"East"`
		North string `yaml:"North"`
		West  string `yaml:"West"`
		South string `yaml:"South"`
	} `yaml:"Exits"`
	Mobiles []struct {
		UUID  string `yaml:"UUID"`
		Count int    `yaml:"Count"`
	} `yaml:"Mobiles"`
}

type YAMLArea struct {
	Items   []YAMLItem   `yaml:"Items"`
	Mobiles []YAMLMobile `yaml:"Mobiles"`
	Rooms   []YAMLRoom   `yaml:"Rooms"`
}

type YAMLPlayer struct {
	Players []YAMLMobile `yaml:"Players"`
}

func buildArea(data []byte) YAMLArea {
	area := YAMLArea{}
	err := yaml.Unmarshal(data, &area)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, item := range area.Items {
		validateItem(item)
	}
	for _, room := range area.Rooms {
		validateRoom(room)
	}
	for _, mobile := range area.Mobiles {
		validateMobile(mobile)
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

func validateMobile(mobile YAMLMobile) {
	if mobile.UUID == "" {
		log.Fatalf("Mobile has no UUID %v", mobile)
	}
	if mobile.Name == "" {
		log.Fatalf("Mobile has no name %v", mobile)
	}
	if mobile.Gear.MainHand == "" {
		log.Fatalf("Mobile gear main hand has no UUID %v", mobile)
	}
}

func validateRoom(room YAMLRoom) {
	if room.UUID == "" {
		log.Fatalf("Room has no UUID %v", room)
	}
	if room.Name == "" {
		log.Fatalf("Room has no Name %v", room)
	}
	if room.Description == "" {
		log.Fatalf("Room has no Description %v", room)
	}
	if room.Size == 0 {
		log.Fatalf("Room size is 0 %v", room)
	}
	for _, mobile := range room.Mobiles {
		if mobile.UUID == "" {
			log.Fatalf("Room mobile has no UUID %v", room)
		}
		if mobile.Count == 0 {
			log.Fatalf("Room mobile has count 0 %v", room)
		}
	}
}

func buildAreaFromPath(path string) YAMLArea {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Area path error")
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("Error reading Area YAML file")
	}
	return buildArea(data)
}

func buildItems(w *world.World, area YAMLArea) {
	for _, r := range area.Items {
		var i item.Item
		if r.Type == "Armor" {
			i = item.NewArmor(r.UUID, r.Name, r.Slot, r.Keywords)
		} else if r.Type == "Weapon" {
			w := item.NewWeapon(r.UUID, r.Name, r.Keywords, r.DamageType)
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

func buildRooms(w *world.World, area YAMLArea) {
	for _, rr := range area.Rooms {
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

func buildPlayers(w *world.World) {
	data := buildPlayerFromPath("data/player.yaml")
	players := YAMLPlayer{}
	err := yaml.Unmarshal(data, &players)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, rp := range players.Players {
		validateMobile(rp)
		c := world.NewPlayer(rp.UUID, rp.Name)

		c.Essence = rp.Essence
		c.Core.Power.Base = rp.Power
		c.Core.Agility.Base = rp.Agility
		c.Core.Insight.Base = rp.Insight
		c.Core.Will.Base = rp.Will

		c.Health.Current = rp.Health
		c.Spirit.Current = rp.Spirit

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

func buildMobiles(w *world.World, area YAMLArea) {
	for _, rp := range area.Mobiles {
		c := world.NewPlayer(rp.UUID, rp.Name)

		c.Essence = rp.Essence
		c.Core.Power.Base = rp.Power
		c.Core.Agility.Base = rp.Agility
		c.Core.Insight.Base = rp.Insight
		c.Core.Will.Base = rp.Will

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

		c.Restore()
		c.Update(0)
		w.Mobiles.AddPrototype(*c)
	}
}

func buildPlayerFromPath(path string) []byte {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Player path error")
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("Error reading player YAML file")
	}
	return data
}

func build() *world.World {
	world := world.NewWorld()
	area := buildAreaFromPath("data/area.yaml")
	buildItems(world, area)
	buildMobiles(world, area)
	buildRooms(world, area)
	buildPlayers(world)
	return world
}
