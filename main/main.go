package main

import (
	"fmt"
	"github.com/RFloTeo/power-spy/display"
	"github.com/RFloTeo/power-spy/resources"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"time"
)

func main() {
	err := resources.InitDocker()
	if err != nil {
		log.Fatalf("Couldn't start Docker client: %s\n", err.Error())
	}

	// Initialise logger
	f, err := os.OpenFile("log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Something happened: %s", err.Error())
		os.Exit(1)
	}
	defer f.Close()
	log.SetOutput(f)

	// Start BubbleTea TUI
	p := tea.NewProgram(initModel())
	if err = p.Start(); err != nil {
		fmt.Printf("Something happened: %s", err.Error())
		os.Exit(1)
	}
}

func initModel() display.Model {
	return display.Model{
		Stats:    resources.Refresh(),
		Duration: 3 * time.Second,
	}
}
