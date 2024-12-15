// backend/internal/services/health_checker.go
package services

import (
	"context"
	"fmt"
	"health-dashboard/internal/models"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trafficmanager/armtrafficmanager"
)

type HealthChecker struct {
	cred           *azidentity.DefaultAzureCredential
	tmClient       *armtrafficmanager.ProfilesClient
	kvClient       *azsecrets.Client
	acrClient      *armcontainerregistry.RegistriesClient
	subscriptionID string
	resourceGroup  string
	region         string
	role           string
	profileName    string // Add this field
}

func NewHealthChecker(subscriptionID, resourceGroup, region, role string) (*HealthChecker, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %v", err)
	}

	tmClient, err := armtrafficmanager.NewProfilesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create traffic manager client: %v", err)
	}

	// TODO: Initialize other clients (kvClient, acrClient)

	return &HealthChecker{
		cred:           cred,
		tmClient:       tmClient,
		subscriptionID: subscriptionID,
		resourceGroup:  resourceGroup,
		region:         region,
		role:           role,
		profileName:    "your-tm-profile-name", // Set this based on your configuration
	}, nil
}

// Add the TriggerFailover method
func (hc *HealthChecker) TriggerFailover(ctx context.Context, targetRegion string) error {
	result, err := hc.tmClient.Get(ctx, hc.resourceGroup, hc.profileName, nil)
	if err != nil {
		return fmt.Errorf("failed to get traffic manager profile: %v", err)
	}

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
		return fmt.Errorf("failed to update traffic manager profile: %v", err)
	}

	return nil
}

// Your existing methods...
func (hc *HealthChecker) CheckHealth(ctx context.Context) (*models.HealthStatus, error) {
	tmStatus := hc.checkTrafficManager(ctx)
	kvStatus := hc.checkKeyVault(ctx)
	acrStatus := hc.checkContainerRegistry(ctx)

	status := &models.HealthStatus{
		TrafficManager:    tmStatus,
		KeyVault:          kvStatus,
		ContainerRegistry: acrStatus,
		RegionStatus:      "healthy", // Determined by overall health
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
	// TODO: Implement actual Key Vault health check
	return "connected"
}

func (hc *HealthChecker) checkContainerRegistry(ctx context.Context) string {
	// TODO: Implement actual Container Registry health check
	return "connected"
}

func (hc *HealthChecker) GetFailoverHistory(ctx context.Context) (*models.FailoverHistory, error) {
	// TODO: Implement actual failover history tracking
	return &models.FailoverHistory{
		LastFailover:   time.Now(),
		CurrentPrimary: hc.region,
		FailoverCount:  1,
	}, nil
}
