// Package memory implements a structure to maintain a running list of
// observed game events. These events are remembered for a finite amount
// of time.  This structure is intended to allow logic elsewhere to
// dynamically adjust based on what's happened recently.
package memory

type GameEvent struct {
	name string
	time int
}

type Memory struct {
	GameEvents []GameEvent
}

func NewMemory() *Memory {
	return &Memory{
		GameEvents: make([]GameEvent, 0),
	}
}

func (m *Memory) AddGameEvent(name string, time int) {
	ge := GameEvent{
		name: name,
		time: time,
	}
	m.GameEvents = append(m.GameEvents, ge)
}

func (m *Memory) MostRecent() string {
	if len(m.GameEvents) == 0 {
		return ""
	}
	return m.GameEvents[0].name
}

func (m *Memory) LastSeen(name string) int {
	total := 0
	for i := len(m.GameEvents) - 1; i >= 0; i-- {
		if m.GameEvents[i].name == name {
			return total
		}
	}
	return -1
}

func (m *Memory) Occurrences(name string) int {
	total := 0
	for _, ge := range m.GameEvents {
		if ge.name == name {
			total += 1
		}
	}
	return total
}

func (m *Memory) Update(tick int) {
	if len(m.GameEvents) == 0 {
		return
	}
	working := make([]GameEvent, 0)
	for _, ge := range m.GameEvents {
		ge.time -= 1
		if ge.time > 0 {
			working = append(working, ge)
		}
	}
	m.GameEvents = working
}
