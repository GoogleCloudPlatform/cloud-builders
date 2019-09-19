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
// Package common functionality shared between commands.
package common

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/deployer"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

// CreateMapFromEqualDelimitedStrings creates a map[string]string from a slice
// of "="-delimited strings.
func CreateMapFromEqualDelimitedStrings(labels []string) (map[string]string, error) {
	labelsMap := make(map[string]string)
	for _, label := range labels {
		p := strings.TrimSpace(label)
		p = strings.Trim(p, ",")
		if p == "" {
			continue
		}
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid key value pair: %q", p)
		}
		k := strings.TrimSpace(kv[0])
		if k == "" {
			return nil, fmt.Errorf("key must not be empty string")
		}
		v := strings.TrimSpace(kv[1])
		if v == "" {
			return nil, fmt.Errorf("value must not be empty string")
		}
		labelsMap[k] = v
	}
	return labelsMap, nil
}

// CreateDeployer creates a Deployer with initialized clients.
func CreateDeployer(ctx context.Context, useGcloud, verbose bool) (*deployer.Deployer, error) {
	c, err := services.NewClients(ctx, useGcloud, verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Clients: %v", err)
	}
	d := &deployer.Deployer{
		Clients:   c,
		UseGcloud: useGcloud,
	}
	return d, nil
}

// SuggestedOutputPath takes a root output directory and returns the path where
// suggested configs should be stored.
func SuggestedOutputPath(root string) string {
	return filepath.Join(root, "suggested")
}

// ExpandedOutputPath takes a root output directory and returns the path where
// expanded configs should be stored.
func ExpandedOutputPath(root string) string {
	return filepath.Join(root, "expanded")
}

// GcloudInPath returns true if the `gcloud` command is in this machine's PATH.
func GcloudInPath() bool {
	if _, err := exec.LookPath("gcloud"); err != nil {
		return false
	}
	return true
}
