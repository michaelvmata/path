package actions

import (
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/items"
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

type GetItemer interface {
	GetItem(string) (item.Item, bool)
}

type Player interface {
	AdjustEssence(int)
	Receive(item.Item) error
	Showln(string, ...interface{})
}

func (qod QuestOnDeath) AssignRewards(World GetItemer, player Player, q *quest.Quest) {
	if q.Reward.Essence > 0 {
		player.Showln("You earned %d essence.", q.Reward.Essence)
		player.AdjustEssence(q.Reward.Essence)
	}
	for _, i := range q.Reward.Items {
		item, ok := World.GetItem(i.UUID)
		if !ok {
			log.Fatalf("Item reward not found %s for quest %s", i.UUID, q.UUID)
		}
		player.Receive(item)
		player.Showln("You earned %s.", item.Name())
	}
}
