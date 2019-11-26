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
// Package cluster contains logic related to talking to a GKE cluster.
package cluster

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/resource"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

// AuthorizeAccess authorizes kubectl to the cluster. In doing so, this also verifies the cluster
// exists.
func AuthorizeAccess(ctx context.Context, clusterName, clusterLocation, clusterProject string, gs services.GcloudService) error {
	if err := gs.ContainerClustersGetCredentials(ctx, clusterName, clusterLocation, clusterProject); err != nil {
		return fmt.Errorf("failed to authorize access: %v", err)
	}
	return nil
}

// ApplyConfigFromString applies a config string to the current context's cluster.
func ApplyConfigFromString(configString, namespace string, ks services.KubectlService) error {
	if err := ks.ApplyFromString(configString, namespace); err != nil {
		return fmt.Errorf("failed to apply config from string: %v", err)
	}
	return nil
}

// GetDeployedObject gets an object deployed to the current context's cluster.
func GetDeployedObject(ctx context.Context, kind, name, namespace string, ks services.KubectlService) (*resource.Object, error) {
	objYaml, err := ks.Get(ctx, kind, name, namespace, "yaml", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get config of deployed object: %v", err)
	}
	return resource.DecodeFromYAML(ctx, []byte(objYaml))
}

// DeployedObjectExists returns true if a deployed object exists in the current context's cluster,
// else false.
func DeployedObjectExists(ctx context.Context, kind, name, namespace string, ks services.KubectlService) (bool, error) {
	objYaml, err := ks.Get(ctx, kind, name, namespace, "yaml", true)
	if err != nil {
		return false, fmt.Errorf("failed to get config of deployed object: %v", err)
	}
	if objYaml == "" {
		return false, nil
	}
	return true, nil
}
