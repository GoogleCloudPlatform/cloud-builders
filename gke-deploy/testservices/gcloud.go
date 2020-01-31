package testservices

import (
	"context"
)

// TestGcloud implements the GcloudService interface.
type TestGcloud struct {
	ContainerClustersGetCredentialsErr error

	ConfigGetValueResp string
	ConfigGetValueErr  error
}

// ContainerClustersGetCredentials calls `gcloud container clusters get-credentials <clusterName> --zone=<clusterLocation> --project=<clusterProject>`.
func (g *TestGcloud) ContainerClustersGetCredentials(ctx context.Context, clusterName, clusterLocation, clusterProject string) error {
	return g.ContainerClustersGetCredentialsErr
}

// ConfigGetValue calls `gcloud config get-value <property>` and returns stdout.
func (g *TestGcloud) ConfigGetValue(ctx context.Context, property string) (string, error) {
	return g.ConfigGetValueResp, g.ConfigGetValueErr
}
