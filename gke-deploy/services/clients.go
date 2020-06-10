// Package services contains logic related to HTTP and CLI clients.
package services

import (
	"context"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Clients is a wrapper around HTTP clients and CLIs.
type Clients struct {
	Gcloud  GcloudService
	Kubectl KubectlService
	OS      OSService
	Remote  RemoteService
	GCS     GcsService
}

// OSService is an interface for os operations.
type OSService interface {
	Stat(ctx context.Context, filename string) (os.FileInfo, error)
	ReadDir(ctx context.Context, dirname string) ([]os.FileInfo, error)
	ReadFile(ctx context.Context, filename string) ([]byte, error)
	WriteFile(ctx context.Context, filename string, data []byte, perm os.FileMode) error
	MkdirAll(ctx context.Context, dirname string, perm os.FileMode) error
	RemoveAll(ctx context.Context, dir string) error
	TempDir(ctx context.Context, dir, pattern string) (string, error)
}

// GcloudService is an interface for gcloud operations.
type GcloudService interface {
	ContainerClustersGetCredentials(ctx context.Context, clusterName, clusterLocation, clusterProject string) error
	ConfigGetValue(ctx context.Context, property string) (string, error)
}

// KubectlService is an interface for kubectl operations.
type KubectlService interface {
	Apply(ctx context.Context, filename, namespace string) error
	ApplyFromString(ctx context.Context, configString, namespace string) error
	Get(ctx context.Context, kind, name, namespace, format string, ignoreNotFound bool) (string, error)
}

// RemoteService is an interface for github.com/google/go-containerregistry/pkg/v1/remote.
type RemoteService interface {
	Image(ref name.Reference) (v1.Image, error)
}

// GcsService is an interface for interacting with Google Cloud Storage
type GcsService interface {
	Copy(ctx context.Context, src, dst string, recursive bool) error
}

// NewClients returns a new Clients object with default services.
func NewClients(ctx context.Context, useGcloud, printCommands bool, serverDryRun bool) (*Clients, error) {
	oss, err := NewOS(ctx)
	if err != nil {
		return nil, err
	}
	var gs GcloudService
	if useGcloud {
		svc, err := NewGcloud(ctx, printCommands)
		if err != nil {
			return nil, err
		}
		gs = svc
	}
	ks, err := NewKubectl(ctx, printCommands, serverDryRun)
	if err != nil {
		return nil, err
	}
	rs, err := NewRemote(ctx)
	if err != nil {
		return nil, err
	}

	ss, err := NewGsutil(ctx, printCommands)
	if err != nil {
		return nil, err
	}

	return &Clients{
		OS:      oss,
		Gcloud:  gs,
		Kubectl: ks,
		Remote:  rs,
		GCS:     ss,
	}, nil
}
