package main

import (
	"fmt"
	"strings"
)

const (
	ESC = "\x1b[38;5;" // ESC character in hex, followed by color code for 256 colors
	END = "m"          // Finishing symbol
)

var Colors = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
}

func Colorize(content string) string {
	args := make([]string, 0)
	for i, color := range Colors {
		args = append(args, fmt.Sprintf("<%s>", color))
		args = append(args, fmt.Sprintf("%s%d%s", ESC, i, END))
	}
	r := strings.NewReplacer(args...)
	return r.Replace(content)
}
