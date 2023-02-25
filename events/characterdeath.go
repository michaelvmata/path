package events

import (
	"github.com/michaelvmata/path/world"
)

type CharacterDeathPayload struct {
	Character *world.Character
	Killer    *world.Character
	World     *world.World
}

type listenerSignature interface {
	Handle(CharacterDeathPayload)
}

type characterDeath struct {
	listeners []listenerSignature
}

var CharacterDeath characterDeath

func (cd *characterDeath) Register(listener listenerSignature) {
	cd.listeners = append(cd.listeners, listener)
}

func (cd *characterDeath) Emit(payload CharacterDeathPayload) {
	for _, listener := range cd.listeners {
		listener.Handle(payload)
	}
}
