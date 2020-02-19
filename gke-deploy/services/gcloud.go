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
	if _, err := runCommand(ctx, g.printCommands, "gcloud", "container", "clusters", "get-credentials", clusterName, fmt.Sprintf("--zone=%s", clusterLocation), fmt.Sprintf("--project=%s", clusterProject), "--quiet"); err != nil {
		return fmt.Errorf("command to get cluster credentials failed: %v", err)
	}
	return nil
}

// ConfigGetValue calls `gcloud config get-value <property>` and returns stdout.
func (g *Gcloud) ConfigGetValue(ctx context.Context, property string) (string, error) {
	out, err := runCommand(ctx, g.printCommands, "gcloud", "config", "get-value", property, "--quiet")
	if err != nil {
		return "", fmt.Errorf("command to get property value failed: %v", err)
	}
	return strings.TrimSpace(out), nil
}
