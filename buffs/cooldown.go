package buffs

type CoolDown struct {
	Lifetime int
	name     string
}

func (cd *CoolDown) Name() string {
	return cd.name
}

func (cd *CoolDown) Update(tick int) {
	cd.Lifetime--
}

func (cd *CoolDown) Expire() {
	cd.Lifetime = 0
}

func (cd *CoolDown) IsExpired() bool {
	return cd.Lifetime <= 0
}

func (cd *CoolDown) Remaining() int {
	return cd.Lifetime
}

func NewCoolDown(lifetime int, name string) CoolDown {
	return CoolDown{
		Lifetime: lifetime,
		name:     name,
	}
}
