package memory

import "testing"

func TestNewMemory(t *testing.T) {
	memory := NewMemory()
	if memory.MostRecent() != "" {
		t.Fatalf("New memory is not empty")
	}
	testEvent := "test event"

	if memory.LastSeen(testEvent) != -1 {
		t.Fatalf("Event unexpectedly seen")
	}
	memory.AddGameEvent(testEvent, 2)
	if memory.MostRecent() != testEvent {
		t.Fatalf("Unexpected most recent event")
	}
	if memory.Occurrences(testEvent) != 1 {
		t.Fatalf("Unexpected event count")
	}

	// Test decay when game event is still remembered
	memory.Update(0)
	if memory.LastSeen(testEvent) != 0 {
		t.Fatalf("Last seen expected=0, actual=%d", memory.LastSeen(testEvent))
	}

	// Test decay when game time is forgotten
	memory.Update(0)
	if memory.LastSeen(testEvent) != -1 {
		t.Fatalf("Last seen expected=-1, actual=%d", memory.LastSeen(testEvent))
	}

	// Test when there are game events
	memory.Update(0)
}
