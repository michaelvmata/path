package simulate

import (
	"fmt"
	"github.com/michaelvmata/path/buffs"
	"github.com/michaelvmata/path/world"
)

func Buffs(w *world.World) {
	for _, player := range w.Players {
		DoBuffs(player)
	}
	for _, mobile := range w.Mobiles.Instances {
		DoBuffs(mobile)
	}
}

func DoBuffs(character *world.Character) {
	for _, buff := range character.Buffs {
		switch buff := buff.(type) {
		case *buffs.Bleed:
			DoBleed(buff)
		}
	}
}

func DoBleed(bleed *buffs.Bleed) {
	damageAmount := 1000
	message := world.Message{
		FirstPerson:         bleed.Applier,
		FirstPersonMessage:  fmt.Sprintf("%s takes %d bleed damage from the wound you inflicted.", bleed.Character.Name, damageAmount),
		SecondPerson:        bleed.Character,
		SecondPersonMessage: fmt.Sprintf("You take %d bleed damage.", damageAmount),
		ThirdPersonMessage:  fmt.Sprintf("%s takes %d bleed damage.", bleed.Character.Name, damageAmount),
	}
	bleed.Character.Room.ShowMessage(message)
	DoDamage(bleed.Applier, bleed.Character, damageAmount)
}
