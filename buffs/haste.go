package buffs

var HasteName = "haste"

type Haste struct {
	Lifetime int
}

func (h Haste) Name() string {
	return HasteName
}

func (h *Haste) Update(tick int) {
	h.Lifetime--
}

func (h Haste) IsExpired() bool {
	return h.Lifetime <= 0
}

func (h Haste) ApplyMessage() string {
	return "The world slows perceptibly."
}

func (h Haste) UnapplyMessage() string {
	return "The world speeds up perceptibly."
}

func NewHaste(level int) *Haste {
	return &Haste{Lifetime: level * 60}
}
