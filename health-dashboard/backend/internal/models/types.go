package models

import "time"

type SystemInfo struct {
	Region           string    `json:"region"`
	Hostname         string    `json:"hostname"`
	ContainerVersion string    `json:"containerVersion"`
	Timestamp        time.Time `json:"timestamp"`
}

type HealthStatus struct {
	TrafficManager    string    `json:"trafficManager"`
	KeyVault          string    `json:"keyVault"`
	ContainerRegistry string    `json:"containerRegistry"`
	RegionStatus      string    `json:"regionStatus"`
	Role              string    `json:"role"`
	LastChecked       time.Time `json:"lastChecked"`
}

type FailoverHistory struct {
	LastFailover   time.Time `json:"lastFailover"`
	CurrentPrimary string    `json:"currentPrimary"`
	FailoverCount  int       `json:"failoverCount"`
}

type FailoverRequest struct {
	TargetRegion string `json:"targetRegion" binding:"required"`
}
