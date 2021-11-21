package main

func main() {
	session := NewSession()
	world := build()
	session.player = world.Players["gaigen"]
	prompt := make(chan bool)
	done := make(chan bool)
	go handleInput(session.incoming)
	go handleOutput(session.outgoing, prompt, done)
	prompt <- true
	for {
		text := <-session.incoming
		command := determineCommand(text)
		command.Execute(world, session, text)
		if text == "quit" {
			done <- true
			break
		}
		prompt <- true
	}
}
