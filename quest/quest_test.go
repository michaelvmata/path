package quest

import "testing"

func TestQuest(t *testing.T) {
	questName := "Test Quest"
	q := NewQuest(questName)
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
}
