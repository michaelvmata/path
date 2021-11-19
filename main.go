package main

func main() {
	session := NewSession()
	prompt := make(chan bool)
	done := make(chan bool)
	go handleInput(session.incoming)
	go handleOutput(session.outgoing, prompt, done)
	prompt <- true
	for {
		text := <-session.incoming
		session.outgoing <- text
		if text == "quit" {
			done <- true
			break
		}
		prompt <- true
	}
}
