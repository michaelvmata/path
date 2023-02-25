package actions

import (
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/world"
	"log"
)

type RespawnCharacter struct{}

func (rp RespawnCharacter) Handle(World *world.World, payload events.CharacterDeathPayload) {
	char := payload.Character
	char.Showln("You were defeated by %s.", payload.Killer.Name)
	for _, c := range char.Room.Players {
		if c == payload.Killer {
			c.Showln("Victory!  %s falls at your hand.", char.Name)
		} else if c != char {
			c.Showln("%s defeated %s", payload.Killer.Name, char.Name)
		}
	}
	if World.IsMobile(char) {
		World.Mobiles.Unspawn(char)
		if err := char.Room.Exit(char); err != nil {
			log.Fatalf("Character died without exiting room %v", payload)
		}
		char.Room = nil
	} else {
		char.Restore()
		char.Showln("In a flash, you're made whole again.")
	}

	opponents := make([]*world.Character, 0)
	for _, opponent := range char.Attacking {
		opponents = append(opponents, opponent)
	}
	for _, opponent := range opponents {
		opponent.StopAttacking(payload.Character)
		char.StopAttacking(opponent)
	}
}
