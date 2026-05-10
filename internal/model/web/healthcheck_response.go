package web

type HealthcheckResponse struct {
	Status     string                `json:"status"`
	SystemInfo HealthcheckSystemInfo `json:"system_info"`
}

type HealthcheckSystemInfo struct {
	Environment string `json:"environment"`
	Version     string `json:"version"`
}
