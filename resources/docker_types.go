package resources

type Container struct {
	Id    string
	Image string
}

type NetworkETH struct {
	RxBytes int `json:"rx_bytes"`
	TxBytes int `json:"tx_bytes"`
}

type CpuStats struct {
	CpuUsage struct {
		Total  float64   `json:"total_usage"`
		PerCPU []float64 `json:"percpu_usage"`
	} `json:"cpu_usage"`
	SystemUsage float64 `json:"system_cpu_usage"`
	OnlineCpus  float64 `json:"online_cpus"`
}

type JsonStats struct {
	Network struct {
		Eth0 NetworkETH `json:"eth0"`
		Eth5 NetworkETH `json:"eth5"`
	} `json:"networks"`
	Memory struct {
		Stats struct {
			Cache int `json:"cache"`
		} `json:"stats"`
		Usage int `json:"usage"`
		Limit int `json:"limit"`
	} `json:"memory_stats"`
	Cpu    CpuStats `json:"cpu_stats"`
	Precpu CpuStats `json:"precpu_stats"`
}

type DockerStats struct {
	Memory        int
	MemoryPercent float64
	CPU           float64
	NetworkIn     int
	NetworkOut    int
}
