// services/health_checker.go
package services

import (
	"context"
	"fmt"
	"health-dashboard/backend/internal/models"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trafficmanager/armtrafficmanager"
)

type HealthChecker struct {
	cred            *azidentity.DefaultAzureCredential
	tmClient        *armtrafficmanager.ProfilesClient
	kvClient        *azsecrets.Client
	acrClient       *armcontainerregistry.RegistriesClient
	subscriptionID  string
	resourceGroup   string
	region          string
	role            string
	profileName     string
	failoverHistory []models.FailoverEvent
	historyMutex    sync.RWMutex
}

func NewHealthChecker(cred *azidentity.DefaultAzureCredential, subscriptionID, resourceGroup, region, role, profileName string) (*HealthChecker, error) {
	tmClient, err := armtrafficmanager.NewProfilesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create traffic manager client: %v", err)
	}

	kvClient, err := azsecrets.NewClient("https://your-keyvault.vault.azure.net/", cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create key vault client: %v", err)
	}

	acrClient, err := armcontainerregistry.NewRegistriesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create container registry client: %v", err)
	}

	return &HealthChecker{
		cred:            cred,
		tmClient:        tmClient,
		kvClient:        kvClient,
		acrClient:       acrClient,
		subscriptionID:  subscriptionID,
		resourceGroup:   resourceGroup,
		region:          region,
		role:            role,
		profileName:     profileName,
		failoverHistory: make([]models.FailoverEvent, 0),
	}, nil
}

func (hc *HealthChecker) CheckHealth(ctx context.Context) (*models.HealthStatus, error) {
	tmStatus := hc.checkTrafficManager(ctx)
	kvStatus := hc.checkKeyVault(ctx)
	acrStatus := hc.checkContainerRegistry(ctx)

	// Determine overall health status
	regionStatus := "healthy"
	if tmStatus != "connected" || kvStatus != "connected" || acrStatus != "connected" {
		regionStatus = "degraded"
	}

	status := &models.HealthStatus{
		TrafficManager:    tmStatus,
		KeyVault:          kvStatus,
		ContainerRegistry: acrStatus,
		RegionStatus:      regionStatus,
		Role:              hc.role,
		LastChecked:       time.Now(),
	}

	return status, nil
}

func (hc *HealthChecker) checkTrafficManager(ctx context.Context) string {
	_, err := hc.tmClient.Get(ctx, hc.resourceGroup, hc.profileName, nil)
	if err != nil {
		return "disconnected"
	}
	return "connected"
}

func (hc *HealthChecker) checkKeyVault(ctx context.Context) string {
	_, err := hc.kvClient.GetSecret(ctx, "test-secret", "", nil)
	if err != nil {
		return "disconnected"
	}
	return "connected"
}

func (hc *HealthChecker) checkContainerRegistry(ctx context.Context) string {
	_, err := hc.acrClient.Get(ctx, hc.resourceGroup, "your-acr-name", nil)
	if err != nil {
		return "disconnected"
	}
	return "connected"
}

func (hc *HealthChecker) TriggerFailover(ctx context.Context, targetRegion string) error {
	result, err := hc.tmClient.Get(ctx, hc.resourceGroup, hc.profileName, nil)
	if err != nil {
		return fmt.Errorf("failed to get traffic manager profile: %v", err)
	}

	currentPrimary := hc.getCurrentPrimary(result.Profile)
	profile := result.Profile

	for i := range profile.Properties.Endpoints {
		var priority int64
		if *profile.Properties.Endpoints[i].Properties.Target == targetRegion {
			priority = 1
		} else {
			priority = 2
		}
		profile.Properties.Endpoints[i].Properties.Priority = &priority
	}

	_, err = hc.tmClient.CreateOrUpdate(ctx, hc.resourceGroup, hc.profileName, profile, nil)
	if err != nil {
		hc.recordFailoverEvent(currentPrimary, targetRegion, "failed", err.Error())
		return fmt.Errorf("failed to update traffic manager profile: %v", err)
	}

	hc.recordFailoverEvent(currentPrimary, targetRegion, "successful", "")
	return nil
}

func (hc *HealthChecker) getCurrentPrimary(profile armtrafficmanager.Profile) string {
	for _, endpoint := range profile.Properties.Endpoints {
		if *endpoint.Properties.Priority == 1 {
			return *endpoint.Properties.Target
		}
	}
	return ""
}

func (hc *HealthChecker) recordFailoverEvent(fromRegion, toRegion, status, errorMessage string) {
	hc.historyMutex.Lock()
	defer hc.historyMutex.Unlock()

	event := models.FailoverEvent{
		Timestamp:    time.Now(),
		FromRegion:   fromRegion,
		ToRegion:     toRegion,
		Status:       status,
		ErrorMessage: errorMessage,
	}

	hc.failoverHistory = append(hc.failoverHistory, event)
}

func (hc *HealthChecker) GetFailoverHistory(ctx context.Context) (*models.FailoverHistory, error) {
	hc.historyMutex.RLock()
	defer hc.historyMutex.RUnlock()

	var currentPrimary string
	if len(hc.failoverHistory) > 0 {
		lastSuccessful := hc.getLastSuccessfulFailover()
		if lastSuccessful != nil {
			currentPrimary = lastSuccessful.ToRegion
		}
	}

	if currentPrimary == "" {
		currentPrimary = hc.region // Default to current region if no failover history
	}

	return &models.FailoverHistory{
		LastFailover:   hc.getLastFailoverTime(),
		CurrentPrimary: currentPrimary,
		FailoverCount:  len(hc.failoverHistory),
		Events:         hc.failoverHistory,
	}, nil
}

func (hc *HealthChecker) getLastFailoverTime() time.Time {
	if len(hc.failoverHistory) == 0 {
		return time.Time{}
	}
	return hc.failoverHistory[len(hc.failoverHistory)-1].Timestamp
}

func (hc *HealthChecker) getLastSuccessfulFailover() *models.FailoverEvent {
	for i := len(hc.failoverHistory) - 1; i >= 0; i-- {
		if hc.failoverHistory[i].Status == "successful" {
			return &hc.failoverHistory[i]
		}
	}
	return nil
}
