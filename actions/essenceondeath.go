package actions

import (
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/world"
)

type EssenceOnDeath struct{}

func (e EssenceOnDeath) Handle(World *world.World, payload events.CharacterDeathPayload) {
	amount := payload.Character.Essence
	payload.Character.DebitEssence(amount)
	killer := payload.Killer
	killer.Showln("%d essence flows to you.", amount)
	killer.CreditEssence(amount)
}
