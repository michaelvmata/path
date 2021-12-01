package skills

import (
	"testing"
)

func TestSkill(t *testing.T) {
	skills := NewSkills()
	if _, ok := skills[DAGGER]; !ok {
		t.Fatalf("Dagger (%d) not in (%v)", DAGGER, skills)
	}
	if _, ok := skills[SWORD]; !ok {
		t.Fatalf("Dagger (%d) not in (%v)", SWORD, skills)
	}
	if _, ok := skills[SPEAR]; !ok {
		t.Fatalf("Spear (%d) not in (%v)", SPEAR, skills)
	}
}
