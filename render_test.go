package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	for _, color := range Colors {
		c := Colorize(fmt.Sprintf("<%s>%s", color, color))
		if !strings.HasPrefix(c, ESC) {
			t.Fatalf("Color %s, %s missing prefix", color, c)
		}
		if !strings.Contains(c, END) {
			t.Fatalf("Color %s, %s missing end", color, c)
		}
		if strings.ContainsAny(c, "<>") {
			t.Fatalf("Color %s, %s didn't strip signs", color, c)
		}
	}
}
