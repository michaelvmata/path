package main

import (
	"bufio"
	"fmt"
	"github.com/michaelvmata/path/session"
	"os"
	"strings"
)

func handleInput(incoming chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		incoming <- strings.TrimSpace(text)
	}
}

func handleOutput(session *session.Session, done chan bool) {
	for {
		select {
		case text := <-session.Outgoing:
			fmt.Print(Colorize(text))
		case <-done:
			break
		}
	}
}
