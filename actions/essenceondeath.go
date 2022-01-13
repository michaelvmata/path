package actions

import "github.com/michaelvmata/path/events"

type EssenceOnDeath struct{}

func (e EssenceOnDeath) Handle(payload events.CharacterDeathPayload) {
	amount := 1
	killer := payload.Killer
	killer.Showln("You gained %d essence.", amount)
	killer.CreditEssence(amount)
}
