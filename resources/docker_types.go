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
		Total int `json:"total_usage"`
	} `json:"cpu_usage"`
	SystemUsage int `json:"system_cpu_usage"`
	OnlineCpus  int `json:"online_cpus"`
}

type JsonStats struct {
	Network struct {
		Eth0 NetworkETH `json:"eth0"`
		Eth5 NetworkETH `json:"eth5"`
	} `json:"network"`
	Memory struct {
		Stats struct {
			Cache int `json:"cache"`
		} `json:"stats"`
		Usage    int `json:"usage"`
		MaxUsage int `json:"max_usage"`
	} `json:"memory_stats"`
	Cpu    CpuStats `json:"cpu_stats"`
	Precpu CpuStats `json:"precpu_stats"`
}

type DockerStats struct {
	Memory        int
	MemoryPercent float32
	CPU           float32
	NetworkIn     int
	NetworkOut    int
}
