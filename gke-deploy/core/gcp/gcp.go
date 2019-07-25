/*
Copyright 2019 Google, Inc. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
