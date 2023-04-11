package quest

import "log"

type Step interface {
	Description() string
	Progress() (int, int)
}

type Quest struct {
	UUID        string
	Description string
	Steps       []Step
}

func NewQuest(UUID string, description string) *Quest {
	return &Quest{
		UUID:        UUID,
		Description: description,
		Steps:       make([]Step, 0),
	}
}

func (q *Quest) Clone(playerUUID string) *Quest {
	cloned := NewQuest(q.UUID, q.Description)
	for _, step := range q.Steps {
		switch s := step.(type) {
		case *KillMobiles:
			km := NewKillMobiles(s.Description(), playerUUID, s.mobileUUID, s.total)
			cloned.Steps = append(cloned.Steps, km)
		default:
			log.Fatalf("Unsupported step %v", step)
		}
	}
	return cloned
}

type KillMobiles struct {
	total       int
	current     int
	description string
	mobileUUID  string
	playerUUID  string
}

func NewKillMobiles(description string, playerUUID string, mobileUUID string, total int) *KillMobiles {
	return &KillMobiles{
		description: description,
		mobileUUID:  mobileUUID,
		playerUUID:  playerUUID,
		total:       total,
	}
}

func (km *KillMobiles) Increment(amount int) {
	km.total += amount
	if km.current > km.total {
		km.current = km.total
	}
}

func (km *KillMobiles) Description() string {
	return km.description
}

func (km *KillMobiles) Progress() (int, int) {
	return km.current, km.total
}
