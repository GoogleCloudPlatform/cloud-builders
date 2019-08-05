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
// Package apply contains the logic for `gke-deploy apply` subcommand.
package apply

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/common"
)

const (
	short = "Skip prepare phase and execute apply phase"
	long  = `Apply Kubernetes configuration files. Skip prepare.

- Apply Kubernetes configuration files to the target cluster with the provided namespace.
- Wait for deployed Kubernetes configuration files to be ready before exiting.
`
	example = `  # Apply only.
  gke-deploy apply -f configs -c my-cluster -n my-namespace -c my-cluster -l us-east1-b

  # Execute prepare and apply, with an intermediary step in between (e.g., manually check expanded YAMLs)
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o expanded -n my-namespace
  cat expanded/*
  gke-deploy apply -f expanded -c my-cluster -n my-namespace -c my-cluster -l us-east1-b  # Pass expanded directory to -f`
)

type options struct {
	filename        string
	clusterLocation string
	clusterName     string
	clusterProject  string
	namespace       string
	verbose         bool
	waitTimeout     time.Duration
}

// NewApplyCommand creates the `gke-deploy apply` subcommand.
func NewApplyCommand() *cobra.Command {
	options := &options{}

	cmd := &cobra.Command{
		Use:     "apply",
		Aliases: []string{"a"},
		Short:   short,
		Long:    long,
		Example: example,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return apply(cmd, options)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&options.filename, "filename", "f", "", "Configuration file or directory of configuration files to use to create Kubernetes objects (file or files in directory must end with \".yml\" or \".yaml\").")
	cmd.Flags().StringVarP(&options.clusterLocation, "location", "l", "", "Region/zone of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.clusterName, "cluster", "c", "", "Name of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.clusterProject, "project", "p", "", "Project of GKE cluster to deploy to. If this field is not provided, the current set GCP project is used.")
	cmd.Flags().StringVarP(&options.namespace, "namespace", "n", "default", "Name of GKE cluster to deploy to.")
	cmd.Flags().BoolVarP(&options.verbose, "verbose", "V", false, "Prints underlying commands being called to stdout.")
	cmd.Flags().DurationVarP(&options.waitTimeout, "timeout", "t", 5*time.Minute, "Timeout limit for waiting for Kubernetes objects to finish applying.")

	return cmd
}

func apply(_ *cobra.Command, options *options) error {
	ctx := context.Background()

	if options.filename == "" {
		return fmt.Errorf("required -f|--filename flag is not set")
	}
	if options.namespace == "" {
		return fmt.Errorf("value of -n|--namespace cannot be empty")
	}
	if options.clusterName != "" && options.clusterLocation == "" {
		return fmt.Errorf("you must set -c|--cluster flag because -l|--location flag is set")
	}
	if options.clusterLocation != "" && options.clusterName == "" {
		return fmt.Errorf("you must set -l|--location flag because -c|--cluster flag is set")
	}

	useGcloud := common.GcloudInPath()
	if !useGcloud && options.clusterName != "" && options.clusterLocation != "" {
		return fmt.Errorf("gcloud must be installed and in PATH to use -c|--cluster and -l|--location")
	}

	d, err := common.CreateDeployer(ctx, useGcloud, options.verbose)
	if err != nil {
		return err
	}

	if err := d.Apply(ctx, options.clusterName, options.clusterLocation, options.clusterProject, options.filename, options.namespace, options.waitTimeout); err != nil {
		return fmt.Errorf("failed to apply deployment: %v", err)
	}

	return nil
}
