package main

import (
	"bufio"
	"fmt"
	"github.com/michaelvmata/path/session"
	"github.com/michaelvmata/path/world"
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

func handleOutput(session *session.Session, done chan bool, player *world.Character) {
	for {
		select {
		case text := <-session.Outgoing:
			fmt.Println(Colorize(text))
		case <-done:
			break
		}
	}
}
