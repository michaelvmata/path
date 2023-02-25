package events

import (
	"github.com/michaelvmata/path/world"
)

type CharacterDeathPayload struct {
	Character *world.Character
	Killer    *world.Character
}

type listenerSignature interface {
	Handle(*world.World, CharacterDeathPayload)
}

type characterDeath struct {
	listeners []listenerSignature
	world     *world.World
}

var CharacterDeath characterDeath

func (cd *characterDeath) Init(World *world.World) {
	cd.world = World
}

func (cd *characterDeath) Register(listener listenerSignature) {
	cd.listeners = append(cd.listeners, listener)
}

func (cd *characterDeath) Emit(payload CharacterDeathPayload) {
	for _, listener := range cd.listeners {
		listener.Handle(cd.world, payload)
	}
}
