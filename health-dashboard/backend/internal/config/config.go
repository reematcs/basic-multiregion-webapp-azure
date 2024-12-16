// internal/config/config.go
package config

import (
	"fmt"
	"os"
)

type Config struct {
	SubscriptionID string
	ResourceGroup  string
	Region         string
	Role           string
	ProfileName    string
	LocalMode      bool
}

func LoadConfig() (*Config, error) {
	subscriptionID := os.Getenv("SUBSCRIPTION_ID")
	if subscriptionID == "" {
		return nil, fmt.Errorf("SUBSCRIPTION_ID environment variable is not set")
	}

	resourceGroup := os.Getenv("RESOURCE_GROUP")
	if resourceGroup == "" {
		return nil, fmt.Errorf("RESOURCE_GROUP environment variable is not set")
	}

	region := os.Getenv("REGION")
	if region == "" {
		return nil, fmt.Errorf("REGION environment variable is not set")
	}

	role := os.Getenv("ROLE")
	if role == "" {
		return nil, fmt.Errorf("ROLE environment variable is not set")
	}

	profileName := os.Getenv("TM_PROFILE_NAME")
	if profileName == "" {
		return nil, fmt.Errorf("TM_PROFILE_NAME environment variable is not set")
	}

	return &Config{
		SubscriptionID: subscriptionID,
		ResourceGroup:  resourceGroup,
		Region:         region,
		Role:           role,
		ProfileName:    profileName,
	}, nil
}

func Load() (*Config, error) {
	config := &Config{
		LocalMode: os.Getenv("LOCAL_MODE") == "true",
	}
	// ... rest of your loading logic
	return config, nil
}
