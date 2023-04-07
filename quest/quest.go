package quest

type Step interface {
	Description() string
	Progress() (int, int)
}

type Quest struct {
	Description string
	Steps       []Step
}

func NewQuest(description string) *Quest {
	return &Quest{
		Description: description,
		Steps:       make([]Step, 0),
	}
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
