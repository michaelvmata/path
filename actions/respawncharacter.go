package actions

import (
	"fmt"
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/world"
	"log"
)

type RespawnCharacter struct{}

func (rp RespawnCharacter) Handle(World *world.World, payload events.CharacterDeathPayload) {
	char := payload.Character
	message := world.Message{
		FirstPerson:  payload.Killer,
		SecondPerson: payload.Character,
	}

	if message.FirstPerson != nil {
		message.FirstPersonMessage = fmt.Sprintf("Victory!  %s falls at your hand.", char.Name)
		message.SecondPersonMessage = fmt.Sprintf("You were defeated by %s.", payload.Killer.Name)
		message.ThirdPersonMessage = fmt.Sprintf("%s defeated %s", payload.Killer.Name, char.Name)
	} else {
		message.SecondPersonMessage = "You succumb to your wounds."
		message.ThirdPersonMessage = fmt.Sprintf("%s succumbs to their wounds.", char.Name)
	}

	message.SecondPerson.Room.ShowMessage(message)

	if World.IsMobile(char) {
		World.Mobiles.Unspawn(char)
		if err := char.Room.Exit(char); err != nil {
			log.Fatalf("Character died without exiting room %v", payload)
		}
		char.Room = nil
	} else {
		if char.Room != char.Anchor && char.Anchor != nil {
			if err := char.Room.Exit(char); err != nil {
				log.Fatalf("Player couldn't exit room %s to respawn.", char.Room.Name())
			}
			char.Room = char.Anchor
			if err := char.Room.Enter(char); err != nil {
				log.Fatalf("Player couldn't respawn in %s", char.Room.Name())
			}
		}
		char.Restore()
		message = world.Message{
			FirstPerson:        char,
			FirstPersonMessage: "In a flash, you're made whole again.",
			ThirdPersonMessage: fmt.Sprintf("In a flash, %s is made whole again.", char.Name),
		}
		char.Room.ShowMessage(message)
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
