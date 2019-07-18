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
// Package run contains the logic for `gke-deploy run` subcommand.
package run

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/common"
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
  gke-deploy run -f configs

  # Deploy to GKE cluster that kubectl is currently targeting without supplying any configs. Have gke-deploy generate base configs for your application using an image, app name, and service port.
  gke-deploy run -i nginx -a nginx -x 80`
)

type options struct {
	appName         string
	appVersion      string
	filename        string
	clusterLocation string
	clusterName     string
	clusterProject  string
	image           string
	labels          []string
	namespace       string
	output          string
	exposePort      int
	verbose         bool
	waitTimeout     time.Duration
}

// NewRunCommand creates the `gke-deploy run` subcommand.
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
	cmd.Flags().StringVarP(&options.filename, "filename", "f", "", "Config file or directory of config files to use to create the Kubernetes resources (file or files in directory must end with \".yml\" or \".yaml\"). If this field is not provided, base configs will be created: Deployment with image provided by --image and HorizontalPodAutoscaler. The application's name will be inferred by the image name's suffix.")
	cmd.Flags().StringVarP(&options.clusterLocation, "location", "l", "", "Region/zone of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.clusterName, "cluster", "c", "", "Name of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.clusterProject, "project", "p", "", "Project of GKE cluster to deploy to. If this field is not provided, the current set GCP project is used.")
	cmd.Flags().StringVarP(&options.image, "image", "i", "", "Image to be deployed.")
	cmd.Flags().StringSliceVarP(&options.labels, "label", "L", nil, "Label(s) to add to Kubernetes resources (k1=v1). Labels can be set comma-delimited or as separate flags. If two or more labels with the same key are listed, the last one is used.")
	cmd.Flags().StringVarP(&options.namespace, "namespace", "n", "default", "Namespace of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.output, "output", "o", "./output", "Target directory to store created and hydrated Kubernetes resource configs. Created configs will be stored in \"<output>/created\" and hydrated configs will be stored in \"<output>/hydrated\".")
	cmd.Flags().IntVarP(&options.exposePort, "expose", "x", 0, "Creates a Service resource that connects to a deployed resource using a selector that matches the label with key as 'app' and value of the image name's suffix specified by --image. The port provided will be used to expose the deployed resource (i.e., port and targetPort will be set to the value provided in this flag).")
	cmd.Flags().BoolVarP(&options.verbose, "verbose", "V", false, "Prints underlying commands being called to stdout.")
	cmd.Flags().DurationVarP(&options.waitTimeout, "timeout", "t", 5*time.Minute, "Timeout limit for waiting for resources to finish applying.")

	return cmd
}

func run(_ *cobra.Command, options *options) error {
	ctx := context.Background()

	var im name.Reference
	if options.image != "" {
		ref, err := name.ParseReference(options.image)
		if err != nil {
			return err
		}
		im = ref
	}

	if options.filename == "" && options.image == "" {
		return fmt.Errorf("omitting -f|--filename requires -i|--image to be set")
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

	if options.exposePort < 0 {
		return fmt.Errorf("value of -x|--expose must be > 0")
	}
	if options.exposePort > 0 && options.image == "" {
		return fmt.Errorf("exposing a deployed resource requires -i|--image to be set")
	}

	labelsMap, err := common.CreateLabelsMap(options.labels)
	if err != nil {
		return err
	}
	d, err := common.CreateDeployer(ctx, options.verbose)
	if err != nil {
		return err
	}

	hydratedOutput := common.HydratedOutputPath(options.output)
	if err := d.Prepare(ctx, im, options.appName, options.appVersion, options.filename, common.CreatedOutputPath(options.output), hydratedOutput, options.namespace, labelsMap, options.exposePort); err != nil {
		return fmt.Errorf("failed to prepare deployment: %v", err)
	}
	if err := d.Apply(ctx, options.clusterName, options.clusterLocation, options.clusterProject, hydratedOutput, options.namespace, options.waitTimeout); err != nil {
		return fmt.Errorf("failed to apply deployment: %v", err)
	}

	return nil
}
