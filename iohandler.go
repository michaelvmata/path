package main

import (
	"bufio"
	"fmt"
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

func handleOutput(session *Session, showPrompt chan bool, done chan bool) {
	prompt := NewPrompt(session)
	for {
		select {
		case text := <-session.outgoing:
			fmt.Println(Colorize(text))
		case <-showPrompt:
			fmt.Printf(Colorize(prompt.Render()))
		case <-done:
			break
		}
	}
}
