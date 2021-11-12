package main

func main() {
	incoming := make(chan string)
	outgoing := make(chan string)
	prompt := make(chan bool)
	done := make(chan bool)
	go handleInput(incoming)
	go handleOutput(outgoing, prompt, done)
	prompt <- true
	for {
		text := <-incoming
		outgoing <- text
		if text == "quit" {
			done <- true
			break
		}
		prompt <- true
	}
}
