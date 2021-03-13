package main

func main() {
	go mainLoop()
	select {}
}
