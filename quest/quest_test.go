package quest

import (
	"testing"
)

func TestQuest(t *testing.T) {
	questName := "Test Quest"
	questUUID := "Test Quest UUID"
	q := NewQuest(questUUID, questName)
	if q.Description != questName {
		t.Fatalf("Quest description note set %s", q.Description)
	}
	playerUUID := "test player"
	mobileUUID := "test mobile"
	description := "Kill mobiles step"
	step := NewKillMobiles(description, playerUUID, mobileUUID, 1)
	current, total := step.Progress()
	if current != 0 || total != 1 {
		t.Fatalf("Progress not set correctly")
	}
	q.Steps = append(q.Steps, step)

	newPlayerUUID := "Test Player UUID"
	cloned := q.Clone(newPlayerUUID)
	if cloned.UUID != cloned.UUID {
		t.Fatalf("Unexpected cloned UUID %s", cloned.UUID)
	}
	for i, _ := range cloned.Steps {
		if cloned.Steps[i].Description() != q.Steps[i].Description() {
			t.Fatalf("Step descriptions don't match at %d", i)
		}
	}
}
