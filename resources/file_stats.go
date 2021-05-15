package resources

import (
	"os"
)

var (
	recording_no = 0
	isRecording  = false
	f            *os.File
)

func toggleRecording() error {
	if isRecording {
		err := f.Close()
		if err != nil {
			return err
		}
	} else {
		var err error
		f, err = os.Create(string(recording_no) + ".csv")
		if err != nil {
			return err
		}
		recording_no++
	}
	isRecording = !isRecording
	return nil
}

func recordStats(stats map[string]DockerStats) {
	//TODO
}
