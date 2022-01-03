package events

import "testing"

type TestCharacterDeathStub struct{}

func (t TestCharacterDeathStub) Handle(payload CharacterDeathPayload) {}

func TestCharacterDeath(t *testing.T) {
	CharacterDeath.Register(TestCharacterDeathStub{})
	CharacterDeath.Emit(CharacterDeathPayload{})
}
