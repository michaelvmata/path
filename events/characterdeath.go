package events

import (
	"github.com/michaelvmata/path/world"
	"log"
)

type CharacterDeathPayload struct {
	Character *world.Character
	Killer    *world.Character
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

func (cd characterDeath) Emit(payload CharacterDeathPayload) {
	log.Println("Emitting")
	for _, listener := range cd.listeners {
		listener.Handle(payload)
	}
}
