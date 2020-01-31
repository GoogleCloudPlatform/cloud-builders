// Package gcp contains logic related to Google Cloud Platform.
package gcp

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

// GetAccount gets the GCP project set during this execution.
func GetProject(ctx context.Context, gs services.GcloudService) (string, error) {
	project, err := gs.ConfigGetValue(ctx, "project")
	if err != nil {
		return "", fmt.Errorf("failed to get project: %v", err)
	}
	return project, nil
}

// GetAccount gets the GCP account set during this execution.
func GetAccount(ctx context.Context, gs services.GcloudService) (string, error) {
	account, err := gs.ConfigGetValue(ctx, "account")
	if err != nil {
		return "", fmt.Errorf("failed to get account: %v", err)
	}
	return account, nil
}
