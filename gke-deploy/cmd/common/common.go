// Package common functionality shared between commands.
package common

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/deployer"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
	applicationsv1beta1 "github.com/kubernetes-sigs/application/pkg/apis/app/v1beta1"
)

// CreateApplicationLinksListFromEqualDelimitedStrings creates a []applicationsv1beta1.Link from a slice
// of "="-delimited strings, where the key is set as Description and the value is set as URL.
func CreateApplicationLinksListFromEqualDelimitedStrings(applicationLinks []string) ([]applicationsv1beta1.Link, error) {
	var asList []applicationsv1beta1.Link
	for _, keyValue := range applicationLinks {
		p := strings.TrimSpace(keyValue)
		p = strings.Trim(p, ",")
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("key value pair %q must be separated by a '=' character", p)
		}
		k := strings.TrimSpace(kv[0])
		if k == "" {
			return nil, fmt.Errorf("key must not be empty string")
		}
		v := strings.TrimSpace(kv[1])
		if v == "" {
			return nil, fmt.Errorf("value must not be empty string")
		}
		asList = append(asList, applicationsv1beta1.Link{
			Description: k,
			URL:         v,
		})
	}
	return asList, nil
}

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
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("key value pair %q must be separated by a '=' character", p)
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
func CreateDeployer(ctx context.Context, useGcloud, verbose bool, serverDryRun bool) (*deployer.Deployer, error) {
	c, err := services.NewClients(ctx, useGcloud, verbose, serverDryRun)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Clients: %v", err)
	}
	d := &deployer.Deployer{
		Clients:      c,
		UseGcloud:    useGcloud,
		ServerDryRun: serverDryRun,
	}
	return d, nil
}

// SuggestedOutputPath takes a root output directory and returns the path where
// suggested configs should be stored.
func SuggestedOutputPath(root string) string {
	return join(root, "suggested")
}

// ExpandedOutputPath takes a root output directory and returns the path where
// expanded configs should be stored.
func ExpandedOutputPath(root string) string {
	return join(root, "expanded")
}

// GcloudInPath returns true if the `gcloud` command is in this machine's PATH.
func GcloudInPath() bool {
	if _, err := exec.LookPath("gcloud"); err != nil {
		return false
	}
	return true
}

func join(base, path string) string {
	u, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}
	u.Path = filepath.Join(u.Path, path)
	return u.String()
}
