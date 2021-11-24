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
			fmt.Println(Colorize(text))
		case <-prompt:
			fmt.Printf(Colorize(" <red>98%s <green>117%s <yellow>85%s <blue>72%s <grey_62>>> "),
				HEART, FIVE_STAR, TWELVE_STAR, CIRCLED_BULLET)
		case <-done:
			break
		}
	}
}
