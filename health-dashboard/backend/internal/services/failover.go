package services

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trafficmanager/armtrafficmanager"
	// Add this import
)

type FailoverService struct {
	tmClient      *armtrafficmanager.ProfilesClient
	resourceGroup string
	profileName   string
}

func (fs *FailoverService) TriggerFailover(ctx context.Context, targetRegion string) error {
	result, err := fs.tmClient.Get(ctx, fs.resourceGroup, fs.profileName, nil)
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

	_, err = fs.tmClient.CreateOrUpdate(ctx, fs.resourceGroup, fs.profileName, profile, nil)
	if err != nil {
		return fmt.Errorf("failed to update traffic manager profile: %v", err)
	}

	return nil
}
