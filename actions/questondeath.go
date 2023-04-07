package actions

import (
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/world"
	"log"
)

func Handle(World *world.World, payload events.CharacterDeathPayload) {
	for _, q := range payload.Killer.Quests {
		log.Printf(q.Description)
	}
}
