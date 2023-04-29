package actions

import (
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/quest"
	"github.com/michaelvmata/path/world"
	"log"
)

type QuestOnDeath struct{}

func (qod QuestOnDeath) Handle(World *world.World, payload events.CharacterDeathPayload) {
	log.Printf("Considering quest on death")
	for _, q := range payload.Killer.Quests {
		for _, step := range q.Steps {
			if killMobiles, ok := step.(*quest.KillMobiles); ok {
				killMobiles.Increment(payload.Killer.UUID, payload.Character.UUID, 1)
			}
		}
		if q.IsComplete() {
			qod.AssignRewards(World, payload.Killer, q)
		}
	}
}

func (qod QuestOnDeath) AssignRewards(World *world.World, player *world.Character, q *quest.Quest) {
	if q.Reward.Essence > 0 {
		player.Showln("You earned %d essence.", q.Reward.Essence)
		player.Essence += q.Reward.Essence
	}
	for _, i := range q.Reward.Items {
		item, ok := World.Items[i.UUID]
		if !ok {
			log.Fatalf("Item reward not found %s for quest %s", i.UUID, q.UUID)
		}
		player.Receive(item)
		player.Showln("You earned %s.", item.Name())
	}
}
