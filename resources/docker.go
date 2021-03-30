package resources

import (
	"net/http"
)

var (
	dClient *http.Client
	dAddr   string
	dCont   string
)

type DockerStats struct {
	CPU        int
	NetworkIn  int
	NetworkOut int
}

func InitDocker(address, container string) {
	dClient = http.DefaultClient
	dAddr = address   // TODO: try figuring out how to get it automagically
	dCont = container // TODO: maybe pass in at GetStats?
}

func GetStats(containers []string) map[string]DockerStats {
	return map[string]DockerStats{"aaa": {1, 2, 3}, "bbb": {7, 1, 56}}
}

func Refresh() map[string]DockerStats {
	return map[string]DockerStats{"aaa": {3, 2, 1}, "bbb": {67, 90, 3}}
}
