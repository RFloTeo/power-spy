package resources

import (
	"log"
	"os"
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
		f, err = os.Create(string(recordingNo) + ".csv")
		if err != nil {
			log.Printf("Failed to start recording: %s\n", err.Error())
			return err
		}
		recordingNo++
	}
	IsRecording = !IsRecording
	return nil
}

func StopRecording() error {
	if IsRecording {
		err := ToggleRecording()
		return err
	}
	return nil
}

func recordStats(stats map[string]DockerStats) {
	//TODO
}
