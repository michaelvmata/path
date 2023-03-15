package help

import "testing"

func TestBuild(t *testing.T) {
	helpIndex := Build("../data/help")
	if len(helpIndex) == 0 {
		t.Fatalf("Failed to load any help files")
	}
}
