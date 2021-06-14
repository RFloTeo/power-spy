package resources

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func ToggleRecording(filter string, containers int) error {
	if IsRecording {
		err := f.Close()
		if err != nil {
			log.Printf("Failed to stop recording: %s\n", err.Error())
			return err
		}
	} else {
		var err error
		f, err = os.Create("recordings/" + filter + strconv.Itoa(recordingNo))
		if err != nil {
			log.Printf("Failed to start recording: %s\n", err.Error())
			return err
		}
		recordingNo++

		startTime = time.Now()
		hasPower := 0
		if PowerOn {
			hasPower = 1
		}
		fmt.Fprintf(f, "%d\n%d\n", containers, hasPower)
	}
	IsRecording = !IsRecording
	return nil
}

var (
	recordingNo = 0
	IsRecording = false
	f           *os.File
	startTime   time.Time
)

// Called upon quitting the program
func StopRecording() error {
	if IsRecording {
		return ToggleRecording("", 0)
	}
	return nil
}

func RecordStats(containers []Container, stats map[string]DockerStats, muW int) {
	if !IsRecording {
		return
	}
	recordTime := time.Since(startTime)
	s := ""
	for _, c := range containers {
		stat := stats[c.Id]
		row := fmt.Sprintf("%s,%s,%d,%f,%f,%d,%d\n", c.Id, c.Image, stat.Memory, stat.MemoryPercent,
			stat.CPU, stat.NetworkIn, stat.NetworkOut)
		s += row
	}
	if PowerOn {
		s += fmt.Sprintln(muW)
	}
	s += fmt.Sprintf("%f\n", recordTime.Seconds())
	n, err := f.WriteString(s)
	if err != nil && len(s) > 0 {
		log.Printf("Tried to record %d bytes of stats, only recorded %d\n", len(s), n)
	}
}
