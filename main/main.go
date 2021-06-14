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

const (
	helpMessage = `A tool for monitoring and recording resource usage stats of Docker containers and power consumption.

Usage: powerspy [COMMANDS]

Commands:
  -d integer  - Set duration between stat fetches in seconds
  -f string   - Set initial filter
  -h          - Display this message and exit
  -l string   - Set location of log file
  -p          - Display and record power readings`
)

var (
	tickTimer  time.Duration
	initFilter string
	logFile    string
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
	wattage := 0
	if resources.PowerOn {
		wattage = resources.GetMicroWatts()
	}
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
		MuW:        wattage,
	}
}

func processFlags() {
	var timer int
	help := flag.Bool("h", false, "Display this message and exit")
	flag.IntVar(&timer, "d", 3, "Set duration between stat fetches in seconds")
	flag.StringVar(&initFilter, "f", "", "Set initial filter")
	flag.StringVar(&logFile, "l", "logs/log.log", "Set location of log file")
	pwr := flag.Bool("p", false, "Display and record power readings")
	flag.Parse()

	tickTimer = time.Duration(timer)
	if *help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}
	resources.PowerOn = *pwr
}
