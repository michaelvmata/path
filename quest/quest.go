package quest

import (
	"fmt"
	"log"
)

type Step interface {
	Description() string
	Progress() (int, int)
	IsComplete() bool
}

type RewardItem struct {
	UUID  string
	Count int
}

type Reward struct {
	Essence int
	Items   []RewardItem
}

func (r *Reward) AddRewardItem(UUID string, count int) {
	r.Items = append(r.Items, RewardItem{
		UUID:  UUID,
		Count: count,
	})
}

type Quest struct {
	UUID        string
	Description string
	Steps       []Step
	Reward      Reward
}

func NewQuest(UUID string, description string) *Quest {
	reward := Reward{
		Essence: 0,
		Items:   make([]RewardItem, 0),
	}
	return &Quest{
		UUID:        UUID,
		Description: description,
		Steps:       make([]Step, 0),
		Reward:      reward,
	}
}

func (q *Quest) IsComplete() bool {
	for _, step := range q.Steps {
		if !step.IsComplete() {
			return false
		}
	}
	return true
}

func (q *Quest) Clone(playerUUID string) *Quest {
	cloned := NewQuest(q.UUID, q.Description)
	for _, step := range q.Steps {
		switch s := step.(type) {
		case *KillMobiles:
			km := NewKillMobiles(s.description, playerUUID, s.mobileUUID, s.total)
			cloned.Steps = append(cloned.Steps, km)
		default:
			log.Fatalf("Unsupported step %v", step)
		}
	}
	cloned.Reward.Essence = q.Reward.Essence
	for _, i := range q.Reward.Items {
		cloned.Reward.AddRewardItem(i.UUID, i.Count)
	}
	return cloned
}

func (q *Quest) Describe() string {
	total := len(q.Steps)
	complete := 0
	for _, s := range q.Steps {
		if s.IsComplete() {
			complete += 1
		}
	}
	return fmt.Sprintf("(%d/%d) %s", complete, total, q.Description)
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

func (km *KillMobiles) Increment(playerUUID string, mobileUUID string, amount int) {
	if playerUUID != km.playerUUID || mobileUUID != km.mobileUUID {
		return
	}
	km.current += amount
	if km.current > km.total {
		km.current = km.total
	}
}

func (km *KillMobiles) Description() string {
	current, total := km.Progress()
	return fmt.Sprintf("(%d/%d) %s", current, total, km.description)
}

func (km *KillMobiles) Progress() (int, int) {
	return km.current, km.total
}

func (km *KillMobiles) IsComplete() bool {
	current, total := km.Progress()
	return current == total
}
