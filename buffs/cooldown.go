package buffs

type CoolDown struct {
	Lifetime int
}

func (cd *CoolDown) Update(tick int) {
	cd.Lifetime--
}

func (cd *CoolDown) Expire() {
	cd.Lifetime = 0
}

func (cd CoolDown) IsExpired() bool {
	return cd.Lifetime <= 0
}
