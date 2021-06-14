package resources

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"math"
)

var (
	dClient *client.Client
)

func InitDocker() error {
	var err error // Necessary because := on the next line creates a local dClient variable
	dClient, err = client.NewClientWithOpts(client.WithVersion("v1.41"))
	return err
}

func GetContainers() ([]Container, error) {
	resp, err := dClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Printf("Getting containers failed: %s\n", err.Error())
		return nil, err
	}
	containers := make([]Container, len(resp))
	for i, c := range resp {
		containers[i] = Container{
			Id:    c.ID,
			Image: c.Image,
		}
	}
	return containers, nil
}

func GetStats(containers []string) map[string]DockerStats {
	m := map[string]DockerStats{}
	for _, c := range containers {
		// Get stats from Docker API
		resp, err := dClient.ContainerStatsOneShot(context.Background(), c)
		if err != nil {
			log.Printf("Fetching stats for container %s failed: %s\n", c, err.Error())
			continue
		}

		// Decode stats from JSON
		decoder := json.NewDecoder(resp.Body)
		var decodedStats JsonStats
		err = decoder.Decode(&decodedStats)
		if err != nil {
			log.Printf("Decoding JSON for container %s failed: %s\n", c, err.Error())
			continue
		}

		// Calculate stats to be displayed and add to map
		usedMemory := decodedStats.Memory.Usage - decodedStats.Memory.Stats.Cache
		cpuDelta := decodedStats.Cpu.CpuUsage.Total - decodedStats.Precpu.CpuUsage.Total
		sysDelta := decodedStats.Cpu.SystemUsage - decodedStats.Precpu.SystemUsage
		cpus := math.Max(decodedStats.Cpu.OnlineCpus, float64(len(decodedStats.Cpu.CpuUsage.PerCPU)))
		stats := DockerStats{
			Memory:        usedMemory,
			MemoryPercent: float64(usedMemory) / float64(decodedStats.Memory.Limit) * 100.0,
			CPU:           (cpuDelta / sysDelta) * cpus * 100.0,
			NetworkIn:     decodedStats.Network.Eth0.RxBytes + decodedStats.Network.Eth5.RxBytes,
			NetworkOut:    decodedStats.Network.Eth0.TxBytes + decodedStats.Network.Eth5.TxBytes,
		}
		m[c] = stats
	}
	return m
}

func Refresh() ([]Container, map[string]DockerStats, error) {
	// Get container list
	cs, err := GetContainers()
	if err != nil {
		log.Println("Refresh failed")
		return []Container{}, nil, err
	}

	// Get stats for all the containers
	ids := make([]string, len(cs))
	for i, c := range cs {
		ids[i] = c.Id
	}
	stats := GetStats(ids)
	return cs, stats, nil
}
