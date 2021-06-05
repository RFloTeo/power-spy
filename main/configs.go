package main

import "time"

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
