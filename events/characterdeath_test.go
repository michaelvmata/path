package events

import (
	"github.com/michaelvmata/path/world"
	"testing"
)

type TestCharacterDeathStub struct{}

func (t TestCharacterDeathStub) Handle(World *world.World, payload CharacterDeathPayload) {}

func TestCharacterDeath(t *testing.T) {
	CharacterDeath.Register(TestCharacterDeathStub{})
	CharacterDeath.Emit(CharacterDeathPayload{})
}
