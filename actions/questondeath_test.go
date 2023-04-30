package actions

import (
	item "github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/quest"
	"testing"
)

type MockPlayer struct{}

func (m MockPlayer) AdjustEssence(amount int)             {}
func (m MockPlayer) Receive(i item.Item) error            { return nil }
func (m MockPlayer) Showln(s string, args ...interface{}) {}

type MockWorld struct{}

func (m MockWorld) GetItem(uuid string) (item.Item, bool) {
	i := item.NewItem("Test Item UUID", "Name", []string{"Test Item"}, "Test Item", "Test Type")
	return i, true
}

func TestQuestOnDeath_AssignRewards(t *testing.T) {
	w := MockWorld{}
	p := MockPlayer{}
	q := quest.NewQuest("Test UUID", "Test Description")
	q.Reward.Essence = 1
	q.Reward.AddRewardItem("Test Item UUID", 1)
	qod := QuestOnDeath{}
	qod.AssignRewards(w, p, q)
}
