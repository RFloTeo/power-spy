package main

import (
	"fmt"
	"github.com/RFloTeo/power-spy/display"
	"github.com/RFloTeo/power-spy/resources"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"time"
)

func main() {
	resources.InitDocker("", "") // TODO: take from command line args?
	p := tea.NewProgram(initModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Something happened: %v", err)
		os.Exit(1)
	}
}

func initModel() display.Model {
	return display.Model{
		Stats:    resources.Refresh(),
		Duration: 3 * time.Second,
	}
}
