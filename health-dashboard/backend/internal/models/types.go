// models/types.go
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

type Metric struct {
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type FailoverHistory struct {
	LastFailover   time.Time       `json:"lastFailover"`
	CurrentPrimary string          `json:"currentPrimary"`
	FailoverCount  int             `json:"failoverCount"`
	Events         []FailoverEvent `json:"events"`
}

type FailoverEvent struct {
	Timestamp    time.Time `json:"timestamp"`
	FromRegion   string    `json:"fromRegion"`
	ToRegion     string    `json:"toRegion"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}
