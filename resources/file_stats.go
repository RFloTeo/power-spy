package resources

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	recordingNo = 0
	IsRecording = false
	f           *os.File
)

func ToggleRecording() error {
	if IsRecording {
		err := f.Close()
		if err != nil {
			log.Printf("Failed to stop recording: %s\n", err.Error())
			return err
		}
	} else {
		var err error
		f, err = os.Create("recordings/" + strconv.Itoa(recordingNo))
		if err != nil {
			log.Printf("Failed to start recording: %s\n", err.Error())
			return err
		}
		recordingNo++
	}
	IsRecording = !IsRecording
	return nil
}

// Called upon quitting the program
func StopRecording() error {
	if IsRecording {
		return ToggleRecording()
	}
	return nil
}

func RecordStats(containers []Container, stats map[string]DockerStats) {
	s := ""
	for _, c := range containers {
		stat := stats[c.Id]
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}
		row := fmt.Sprintf("%s,%s,%d,%f,%f,%d,%d\n", c.Id, name, stat.Memory, stat.MemoryPercent,
			stat.CPU, stat.NetworkIn, stat.NetworkOut)
		s += row
	}
	n, err := f.WriteString(s)
	if err != nil && len(s) > 0 {
		log.Printf("Tried to record %d bytes of stats, only recorded %d\n", len(s), n)
	}
}
