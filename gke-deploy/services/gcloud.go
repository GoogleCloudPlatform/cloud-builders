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
package services

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Gcloud implements the GcloudService interface.
type Gcloud struct {
	printCommands bool
}

// NewGcloud returns a new Gcloud object.
func NewGcloud(ctx context.Context, printCommands bool) (*Gcloud, error) {
	if _, err := exec.LookPath("gcloud"); err != nil {
		return nil, err
	}
	return &Gcloud{
		printCommands,
	}, nil
}

// ContainerClustersGetCredentials calls `gcloud container clusters get-credentials <clusterName> --zone=<clusterLocation> --project=<clusterProject>`.
// Both region and zone can be passed to the --zone flag.
func (g *Gcloud) ContainerClustersGetCredentials(ctx context.Context, clusterName, clusterLocation, clusterProject string) error {
	if _, err := runCommand(g.printCommands, "gcloud", "container", "clusters", "get-credentials", clusterName, fmt.Sprintf("--zone=%s", clusterLocation), fmt.Sprintf("--project=%s", clusterProject), "--quiet"); err != nil {
		return fmt.Errorf("command to get cluster credentials failed: %v", err)
	}
	return nil
}

// ConfigGetValue calls `gcloud config get-value <property>` and returns stdout.
func (g *Gcloud) ConfigGetValue(ctx context.Context, property string) (string, error) {
	out, err := runCommand(g.printCommands, "gcloud", "config", "get-value", property, "--quiet")
	if err != nil {
		return "", fmt.Errorf("command to get property value failed: %v", err)
	}
	return strings.TrimSpace(out), nil
}
