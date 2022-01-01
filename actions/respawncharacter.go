package actions

import (
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/world"
	"log"
)

type RespawnCharacter struct{}

func (rp RespawnCharacter) Handle(payload events.CharacterDeathPayload) {
	char := payload.Character
	payload.World.Mobiles.Unspawn(char)
	if err := char.Room.Exit(char); err != nil {
		log.Fatalf("Character died without exiting room %v", payload)
	}
	char.Room = nil

	opponents := make([]*world.Character, 0)
	for _, opponent := range char.Attacking {
		opponents = append(opponents, opponent)
	}
	for _, opponent := range opponents {
		opponent.StopAttacking(payload.Character)
		char.StopAttacking(opponent)
	}
}
