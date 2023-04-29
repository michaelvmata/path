package main

import (
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/quest"
	"github.com/michaelvmata/path/world"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type YAMLItem struct {
	UUID          string   `yaml:"UUID"`
	Name          string   `yaml:"Name"`
	Type          string   `yaml:"Type"`
	Description   string   `yaml:"Description"`
	Slot          string   `yaml:"Slot"`
	DamageType    string   `yaml:"DamageType"`
	Attributes    []string `yaml:"Attributes""`
	Immovable     bool     `yaml:"Immovable"`
	MinimumDamage int      `yaml:"MinimumDamage"`
	MaximumDamage int      `yaml:"MaximumDamage"`
	CriticalRate  float64  `yaml:"CriticalRate"`
	CriticalBonus float64  `yaml:"CriticalBonus"`
	Modifiers     []struct {
		Type  string `yaml:"Type"`
		Value int    `yaml:"Value"`
	} `yaml:"Modifiers"`
	Keywords []string `yaml:"Keywords"`
}

type YAMLMobile struct {
	UUID         string `yaml:"UUID"`
	Name         string `yaml:"Name"`
	Description  string `yaml:"Description"`
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
	Anchor       string `yaml:"Anchor"`
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
		Barrier  int `yaml:"Barrier"`
		Bash     int `yaml:"Bash"`
		Backstab int `yaml:"Backstab"`
		Bleed    int `yaml:"Bleed"`
		Blitz    int `yaml:"Blitz"`
		Circle   int `yaml:"Circle"`
		Evasion  int `yaml:"Evasion"`
		Haste    int `yaml:"Haste"`
		Parry    int `yaml:"Parry"`
		Sweep    int `yaml:"Sweep"`
	} `yaml:"Skills"`
	Quests []YAMLMobileQuests `yaml:"Quests"`
}

type YAMLMobileQuests struct {
	UUID  string                `yaml:"UUID"`
	Steps []YAMLMobileQuestStep `yaml:"Steps"`
}

type YAMLMobileQuestStep struct {
	Current int `yaml:"Current"`
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
	Items []struct {
		UUID  string `yaml:"UUID"`
		Count int    `yaml:"Count"`
	} `yaml:"Items"`
}

type YamlQuest struct {
	UUID        string `yaml:"UUID"`
	Description string `yaml:"Description"`
	Steps       []struct {
		Type        string `yaml:"Type"`
		Total       int    `yaml:"Total"`
		Description string `yaml:"Description"`
		Mobile      string `yaml:"Mobile"`
	} `yaml:"Steps"`
	Rewards struct {
		Essence int `yaml:"Essence"`
		Items   []struct {
			UUID  string `yaml:"UUID"`
			Count int    `yaml:"Count"`
		} `yaml:"Items"`
	} `yaml:"Rewards"`
}

type YAMLArea struct {
	UUID    string       `yaml:"UUID"`
	Name    string       `yaml:"Name"`
	Items   []YAMLItem   `yaml:"Items"`
	Mobiles []YAMLMobile `yaml:"Mobiles"`
	Rooms   []YAMLRoom   `yaml:"Rooms"`
	Quests  []YamlQuest  `yaml:"Quests"`
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
	if area.Name == "" {
		log.Fatalf("Area has no Name")
	}
	if area.UUID == "" {
		log.Fatalf("Area has no UUID")
	}
	for _, item := range area.Items {
		validateItem(item)
	}
	for _, room := range area.Rooms {
		validateRoom(room)
	}
	for _, quest := range area.Quests {
		validateQuest(quest)
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
	if item.Type == "Armor" && item.Slot == "" {
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

func validateQuest(quest YamlQuest) {
	if quest.UUID == "" {
		log.Fatalf("Quest has no UUID %v", quest)
	}
	if quest.Description == "" {
		log.Fatalf("Quest has no description %s", quest.UUID)
	}
	if len(quest.Steps) == 0 {
		log.Fatalf("No steps for quest %s", quest.UUID)
	}
	if quest.Rewards.Essence == 0 {
		log.Fatalf("No rewards for quest %s", quest.UUID)
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

func buildAreas(w *world.World, root string) {
	nodes, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("Error reading areas directory")
	}

	for _, f := range nodes {
		if f.IsDir() {
			continue
		}
		yamlArea := buildAreaFromPath(root + "/" + f.Name())
		area := world.NewArea(yamlArea.UUID, yamlArea.Name)
		w.Areas[area.UUID] = area
		buildItems(w, yamlArea)
		buildMobiles(w, yamlArea)
		buildRooms(w, yamlArea)
		buildQuests(w, yamlArea)
	}
}

func buildItems(w *world.World, area YAMLArea) {
	for _, r := range area.Items {
		var i item.Item
		if r.Type == item.ArmorType {
			i = item.NewArmor(r.UUID, r.Name, r.Slot, r.Keywords, r.Description)
		} else if r.Type == item.WeaponType {
			w := item.NewWeapon(r.UUID, r.Name, r.Keywords, r.Description, r.DamageType, r.Attributes)
			if r.MaximumDamage <= r.MinimumDamage || r.MinimumDamage <= 0 {
				log.Fatalln("Invalid Maximum and Minimum Damage", r)
			}
			w.MinimumDamage = r.MinimumDamage
			w.MaximumDamage = r.MaximumDamage
			w.CriticalBonus = r.CriticalBonus
			w.CriticalRate = r.CriticalRate
			i = w
		} else {
			i = item.NewItem(r.UUID, r.Name, r.Keywords, r.Description, r.Type)
			if r.Immovable {
				i.MakeImmovable()
			}
		}
		for _, rm := range r.Modifiers {
			i.AddModifier(rm.Type, rm.Value)
		}
		w.Items[i.UUID()] = i
	}
}

func buildRooms(w *world.World, yamlArea YAMLArea) {
	area, found := w.Areas[yamlArea.UUID]
	if !found {
		log.Fatalf("Area index not setup")
	}
	for _, rr := range yamlArea.Rooms {
		room := world.NewRoom(rr.UUID, rr.Name, rr.Description, rr.Size, area)
		area.Rooms[room.UUID] = room
		w.Rooms[room.UUID] = room
		for _, mc := range rr.Mobiles {
			w.RoomMobiles[room.UUID] = append(w.RoomMobiles[room.UUID], world.NewRoomMobile(mc.UUID, mc.Count))
		}
		for _, yamlItem := range rr.Items {
			i, found := w.Items[yamlItem.UUID]
			if !found {
				log.Fatalf("Can't find item %s", yamlItem.UUID)
			}
			room.Accept(i)
		}
		room.Exits.East = rr.Exits.East
		room.Exits.North = rr.Exits.North
		room.Exits.South = rr.Exits.South
		room.Exits.West = rr.Exits.West
	}
}

func buildQuests(w *world.World, yamlArea YAMLArea) {
	for _, yamlQuest := range yamlArea.Quests {
		q := quest.NewQuest(yamlQuest.UUID, yamlQuest.Description)
		for _, yamlStep := range yamlQuest.Steps {
			if yamlStep.Type == "KillMobiles" {
				s := quest.NewKillMobiles(yamlStep.Description, "", yamlStep.Mobile, yamlStep.Total)
				q.Steps = append(q.Steps, s)
			}
		}
		q.Reward.Essence = yamlQuest.Rewards.Essence
		for _, yamlItem := range yamlQuest.Rewards.Items {
			q.Reward.AddRewardItem(yamlItem.UUID, yamlItem.Count)
		}
		w.Quests[q.UUID] = q
	}
}

func savePlayers(players map[string]*world.Character) {
	yamlPlayer := YAMLPlayer{
		Players: make([]YAMLMobile, 0),
	}
	for _, player := range players {
		p := YAMLMobile{Inventory: make([]string, 0)}
		p.Name = player.Name
		p.UUID = player.UUID
		p.Essence = player.Essence
		p.RoomUUID = player.Room.UUID

		p.Power = player.Core.Power.Base
		p.Agility = player.Core.Agility.Base
		p.Insight = player.Core.Insight.Base
		p.Will = player.Core.Will.Base

		p.Health = player.Health.Current
		p.Spirit = player.Spirit.Current

		if player.Gear.Head != nil {
			p.Gear.Head = player.Gear.Head.UUID()
		}
		if player.Gear.Neck != nil {
			p.Gear.Neck = player.Gear.Neck.UUID()
		}
		if player.Gear.Body != nil {
			p.Gear.Body = player.Gear.Body.UUID()
		}
		if player.Gear.Arms != nil {
			p.Gear.Arms = player.Gear.Arms.UUID()
		}
		if player.Gear.Hands != nil {
			p.Gear.Hands = player.Gear.Hands.UUID()
		}
		if player.Gear.Waist != nil {
			p.Gear.Waist = player.Gear.Waist.UUID()
		}
		if player.Gear.Legs != nil {
			p.Gear.Legs = player.Gear.Legs.UUID()
		}
		if player.Gear.Feet != nil {
			p.Gear.Feet = player.Gear.Feet.UUID()
		}
		if player.Gear.Wrist != nil {
			p.Gear.Wrist = player.Gear.Wrist.UUID()
		}
		if player.Gear.Fingers != nil {
			p.Gear.Fingers = player.Gear.Fingers.UUID()
		}
		if player.Gear.MainHand != nil {
			p.Gear.MainHand = player.Gear.MainHand.UUID()
		}
		if player.Gear.OffHand != nil {
			p.Gear.OffHand = player.Gear.OffHand.UUID()
		}

		for _, i := range player.Inventory.Items {
			p.Inventory = append(p.Inventory, i.UUID())
		}
		p.Skills.Barrier = player.Skills.Barrier.Base
		p.Skills.Bash = player.Skills.Bash.Base
		p.Skills.Backstab = player.Skills.Backstab.Base
		p.Skills.Bleed = player.Skills.Bleed.Base
		p.Skills.Blitz = player.Skills.Blitz.Base
		p.Skills.Circle = player.Skills.Circle.Base
		p.Skills.Evasion = player.Skills.Evasion.Base
		p.Skills.Haste = player.Skills.Haste.Base
		p.Skills.Parry = player.Skills.Parry.Base
		p.Skills.Sweep = player.Skills.Sweep.Base

		for _, q := range player.Quests {
			steps := make([]YAMLMobileQuestStep, 0)
			for _, s := range q.Steps {
				current, _ := s.Progress()
				steps = append(steps, YAMLMobileQuestStep{Current: current})
			}
			p.Quests = append(p.Quests, YAMLMobileQuests{
				UUID:  q.UUID,
				Steps: steps,
			})
		}
		yamlPlayer.Players = append(yamlPlayer.Players, p)
	}
	data, err := yaml.Marshal(&yamlPlayer)
	if err != nil {
		log.Fatalf("Unable to marshal players yaml")
	}
	yamlFile, err := os.OpenFile("data/player.yaml", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Unable to open players file")
	}
	if err := yamlFile.Truncate(0); err != nil {
		log.Fatalf("Unable to truncate player file")
	}
	len, err := yamlFile.Write(data)
	if err != nil {
		log.Fatalf("Unable to write players file")
	}
	log.Printf("Wrote %d bytes to players file.", len)
	yamlFile.Close()
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
		c.Skills.Backstab.Base = rp.Skills.Backstab
		c.Skills.Bash.Base = rp.Skills.Bash
		c.Skills.Barrier.Base = rp.Skills.Barrier
		c.Skills.Bleed.Base = rp.Skills.Bleed
		c.Skills.Blitz.Base = rp.Skills.Blitz
		c.Skills.Circle.Base = rp.Skills.Circle
		c.Skills.Evasion.Base = rp.Skills.Evasion
		c.Skills.Haste.Base = rp.Skills.Haste
		c.Skills.Parry.Base = rp.Skills.Parry
		c.Skills.Sweep.Base = rp.Skills.Sweep
		c.Update(0)
		w.Players[c.Name] = c
		if room, ok := w.Rooms[rp.RoomUUID]; ok {
			if err := room.Enter(c); err == nil {
				c.Room = room
			}
		}
		if room, ok := w.Rooms[rp.Anchor]; ok {
			c.Anchor = room
		}
		for _, yq := range rp.Quests {
			q, ok := w.Quests[yq.UUID]
			if !ok {
				log.Fatalf("Could not find quest %s", yq.UUID)
			}
			c.Quests = append(c.Quests, q.Clone(c.UUID))
		}
	}
}

func buildMobiles(w *world.World, area YAMLArea) {
	for _, rp := range area.Mobiles {
		c := world.NewPlayer(rp.UUID, rp.Name)

		c.Description = rp.Description
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

func build(root string) *world.World {
	world := world.NewWorld()
	buildAreas(world, root)
	buildPlayers(world)
	return world
}
