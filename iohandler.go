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

func handleOutput(outgoing chan string, prompt chan bool, done chan bool) {
	for {
		select {
		case text := <-outgoing:
			fmt.Println(text)
		case <-prompt:
			fmt.Print(">> ")
		case <-done:
			break
		}
	}
}
