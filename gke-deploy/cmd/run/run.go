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
package run

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/common"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/image"
)

const (
	short = "Execute both prepare and apply phase"
	long  = `Deploy to GKE in two phases, which will do the following:

Prepare Phase:
  - Modify Kubernetes config YAMLs:
    - Set the digest of images that match the [--image|-i] flag, if provided.
    - Add app.kubernetes.io/name=[--app|-a] label, if provided.
    - Add app.kubernetes.io/version=[--version|-v] label, if provided.

Apply Phase:
  - Apply Kubernetes config YAMLs to the target cluster with the provided namespace.
  - Wait for deployed resources to be ready before exiting.
`
	example = `  # Modify configs and deploy to GKE cluster.
  gke-deploy run -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o modified -n my-namespace -c my-cluster -l us-east1-b

  # Deploy to GKE cluster that kubectl is currently targeting.
  gke-deploy run -f configs`
)

type options struct {
	appName         string
	appVersion      string
	filename        string
	clusterLocation string
	clusterName     string
	clusterProject  string
	images          []string
	labels          []string
	namespace       string
	output          string
	verbose         bool
	waitTimeout     time.Duration
}

func NewRunCommand() *cobra.Command {
	options := &options{}

	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"r", "deploy", "d"},
		Short:   short,
		Long:    long,
		Example: example,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return run(cmd, options)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&options.appName, "app", "a", "", "Application name of the Kubernetes deployment.")
	cmd.Flags().StringVarP(&options.appVersion, "version", "v", "", "Version of the Kubernetes deployment.")
	cmd.Flags().StringVarP(&options.filename, "filename", "f", "", "Config file or directory of config files to use to create the Kubernetes resources (file or files in directory must end with \".yml\" or \".yaml\").")
	cmd.Flags().StringVarP(&options.clusterLocation, "location", "l", "", "Region/zone of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.clusterName, "cluster", "c", "", "Name of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.clusterProject, "project", "p", "", "Project of GKE cluster to deploy to. If this field is not provided, the current set GCP project is used.")
	cmd.Flags().StringSliceVarP(&options.images, "image", "i", nil, "Image(s) to be deployed. Images can be set comma-delimited or as separate flags.")
	cmd.Flags().StringSliceVarP(&options.labels, "label", "L", nil, "Label(s) to add to Kubernetes resources (k1=v1). Labels can be set comma-delimited or as separate flags. If two or more labels with the same key are listed, the last one is used.")
	cmd.Flags().StringVarP(&options.namespace, "namespace", "n", "default", "Name of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.output, "output", "o", "./output", "Target directory to store modified Kubernetes resource configs.")
	cmd.Flags().BoolVarP(&options.verbose, "verbose", "V", false, "Prints underlying commands being called to stdout.")
	cmd.Flags().DurationVarP(&options.waitTimeout, "timeout", "t", 5*time.Minute, "Timeout limit for waiting for resources to finish applying.")

	return cmd
}

func run(cmd *cobra.Command, options *options) error {
	ctx := context.Background()

	images, err := image.ParseReferences(options.images)
	if err != nil {
		return err
	}

	if options.filename == "" {
		// TODO(joonlim): Generate base configs if user does not supply any.
		return fmt.Errorf("required -f|--filename flag is not set")
	}
	if options.namespace == "" {
		return fmt.Errorf("value of -n|--namespace cannot be empty")
	}
	if options.output == "" {
		return fmt.Errorf("value of -o|--output cannot be empty")
	}
	if options.clusterName != "" && options.clusterLocation == "" {
		return fmt.Errorf("you must set -c|--cluster flag because -l|--location flag is set")
	}
	if options.clusterLocation != "" && options.clusterName == "" {
		return fmt.Errorf("you must set -l|--location flag because -c|--cluster flag is set")
	}

	labelsMap, err := common.CreateLabelsMap(options.labels)
	if err != nil {
		return err
	}
	d, err := common.CreateDeployer(ctx, options.verbose)
	if err != nil {
		return err
	}

	if err := d.Prepare(ctx, images, options.appName, options.appVersion, options.filename, options.output, options.namespace, labelsMap); err != nil {
		return fmt.Errorf("failed to prepare deployment: %v", err)
	}
	if err := d.Apply(ctx, options.clusterName, options.clusterLocation, options.clusterProject, options.output, options.namespace, options.waitTimeout); err != nil {
		return fmt.Errorf("failed to apply deployment: %v", err)
	}

	return nil
}
