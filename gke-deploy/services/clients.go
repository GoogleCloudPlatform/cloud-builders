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
	"os"
)

// Clients is a wrapper around HTTP clients and CLIs.
type Clients struct {
	Gcloud  GcloudService
	Kubectl KubectlService
	OS      OSService
}

// OSService is an interface for os operations.
type OSService interface {
	Stat(ctx context.Context, filename string) (os.FileInfo, error)
	ReadDir(ctx context.Context, dirname string) ([]os.FileInfo, error)
	ReadFile(ctx context.Context, filename string) ([]byte, error)
	WriteFile(ctx context.Context, filename string, data []byte, perm os.FileMode) error
	MkdirAll(ctx context.Context, dirname string, perm os.FileMode) error
}

// GcloudService is an interface for gcloud operations.
type GcloudService interface {
	ContainerImagesDescribe(ctx context.Context, image, format string) (string, error)
	ContainerClustersGetCredentials(ctx context.Context, clusterName, clusterLocation, clusterProject string) error
	ConfigGetValue(ctx context.Context, property string) (string, error)
}

// KubectlService is an interface for kubectl operations.
type KubectlService interface {
	Apply(ctx context.Context, configs, namespace string) error
	Get(ctx context.Context, kind, name, namespace, format string) (string, error)
}

// NewClients returns a new Clients object with default services.
func NewClients(ctx context.Context, printCommands bool) (*Clients, error) {
	oss, err := NewOS(ctx)
	if err != nil {
		return nil, err
	}
	gs, err := NewGcloud(ctx, printCommands)
	if err != nil {
		return nil, err
	}
	ks, err := NewKubectl(ctx, printCommands)
	if err != nil {
		return nil, err
	}

	return &Clients{
		OS:      oss,
		Gcloud:  gs,
		Kubectl: ks,
	}, nil
}
