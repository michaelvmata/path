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
	total := 2
	step := NewKillMobiles(description, playerUUID, mobileUUID, total)
	current, total := step.Progress()
	if current != 0 || total != total {
		t.Fatalf("Progress not set correctly")
	}

	step.Increment(playerUUID, "Other Mobile UUID", 1)
	current, total = step.Progress()
	if current != 0 {
		t.Fatalf("Progress incremented for wrong mobile")
	}
	step.Increment(playerUUID, mobileUUID, 1)
	current, total = step.Progress()
	if current != 1 {
		t.Fatalf("Progress not incremented for mobile.  Current = %d", current)
	}

	q.Steps = append(q.Steps, step)

	if current == total && q.IsComplete() {
		t.Fatalf("Quest should not be complete yet")
	}
	step.Increment(playerUUID, mobileUUID, total-current)
	if !q.IsComplete() {
		t.Fatalf("Quest should be complete")
	}

	newPlayerUUID := "Test Player UUID"
	cloned := q.Clone(newPlayerUUID)
	if cloned.UUID != cloned.UUID {
		t.Fatalf("Unexpected cloned UUID %s", cloned.UUID)
	}

	if q.Describe() == "" {
		t.Fatalf("Quest description empty")
	}
}
