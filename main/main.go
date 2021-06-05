package main

import (
	"flag"
	"fmt"
	"github.com/RFloTeo/power-spy/display"
	"github.com/RFloTeo/power-spy/resources"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"time"
)

func main() {
	processFlags()

	err := resources.InitDocker()
	if err != nil {
		log.Fatalf("Couldn't start Docker client: %s\n", err.Error())
	}

	// Initialise logger
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
	ti := textinput.NewModel()
	ti.Placeholder = "Filter"
	ti.Blur()
	ti.Width = 30
	ti.CharLimit = 100

	containers, stats, err := resources.Refresh()
	if err != nil {
		log.Println("Failed initial refresh")
	}
	return display.Model{
		Containers: containers,
		Stats:      stats,
		Duration:   tickTimer * time.Second,
		ToggleFail: false,
		StopFail:   false,
		Text:       ti,
		Filter:     initFilter,
	}
}

func processFlags() {
	help := flag.Bool("h", false, "Display this message and exit")
	flag.DurationVar(&tickTimer, "d", 3, "Set duration between stat fetches in seconds")
	flag.StringVar(&initFilter, "f", "", "Set initial filter")
	flag.StringVar(&logFile, "l", "logs/log.log", "Set location of log file")
	pwr := flag.Bool("p", false, "Display and record power readings")
	flag.Parse()

	if *help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}
	resources.PowerOn = *pwr
}
