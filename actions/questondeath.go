package actions

import (
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/quest"
	"github.com/michaelvmata/path/world"
	"log"
)

type QuestOnDeath struct{}

func (q QuestOnDeath) Handle(World *world.World, payload events.CharacterDeathPayload) {
	log.Printf("Considering quest on death")
	for _, q := range payload.Killer.Quests {
		for _, step := range q.Steps {
			if killMobiles, ok := step.(*quest.KillMobiles); ok {
				killMobiles.Increment(payload.Killer.UUID, payload.Character.UUID, 1)
			}
		}
	}
}
