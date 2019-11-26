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
// Package prepare contains the logic for `gke-deploy prepare` subcommand.
package prepare

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/common"
)

const (
	short = "Execute prepare phase and skip apply phase"
	long  = `Prepare to deploy to GKE by generating expanded Kubernetes configuration files. Skip apply.

- Expand Kubernetes configuration files to:
  - Set the digest of images that match the [--image|-i] flag, if provided.
  - Add app.kubernetes.io/name=[--app|-a] label, if provided.
  - Add app.kubernetes.io/version=[--version|-v] label, if provided.
`
	example = `  # Prepare only.
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o expanded -n my-namespace

  # Execute prepare and apply, with an intermediary step in between (e.g., manually check expanded YAMLs)
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o expanded -n my-namespace
  cat expanded/*
  gke-deploy apply -f expanded -c my-cluster -n my-namespace -c my-cluster -l us-east1-b  # Pass expanded directory to -f

  # Pipe output from another templating engine to gke-deploy prepare.
  kustomize build overlays/staging | gke-deploy prepare -f - -a my-app
  helm template charts/prometheus | gke-deploy prepare -f - -a prometheus`
)

type options struct {
	appName     string
	appVersion  string
	filename    string
	image       string
	labels      []string
	annotations []string
	namespace   string
	output      string
	exposePort  int
	verbose     bool
	recursive   bool
}

// NewPrepareCommand creates the `gke-deploy prepare` subcommand.
func NewPrepareCommand() *cobra.Command {
	options := &options{}

	cmd := &cobra.Command{
		Use:     "prepare",
		Aliases: []string{"p"},
		Short:   short,
		Long:    long,
		Example: example,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return prepare(cmd, options)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&options.appName, "app", "a", "", "Application name of the Kubernetes deployment.")
	cmd.Flags().StringVarP(&options.appVersion, "version", "v", "", "Version of the Kubernetes deployment.")
	cmd.Flags().StringVarP(&options.filename, "filename", "f", "", "Configuration file or directory of configuration files to use to create Kubernetes objects (file or files in directory must end with \".yml\" or \".yaml\"). If this field is not provided, suggested base configs will be created: Deployment with image provided by --image and HorizontalPodAutoscaler. The application's name will be inferred by the image name's suffix.")
	cmd.Flags().StringVarP(&options.image, "image", "i", "", "Image to be deployed.")
	cmd.Flags().StringSliceVarP(&options.labels, "label", "L", nil, "Label(s) to add to Kubernetes objects (k1=v1). Labels can be set comma-delimited or as separate flags. If two or more labels with the same key are listed, the last one is used.")
	cmd.Flags().StringVarP(&options.namespace, "namespace", "n", "", "Namespace of GKE cluster to deploy to. Creates a namespace Kubernetes configuration file to reflect this and updates the namespace field of each supplied Kubernetes configuration file.")
	cmd.Flags().StringSliceVarP(&options.annotations, "annotation", "A", nil, "Annotation(s) to add to Kubernetes configuration files (k1=v1). Annotations can be set comma-delimited or as separate flags. If two or more annotations with the same key are listed, the last one is used.")
	cmd.Flags().StringVarP(&options.output, "output", "o", "./output", "Target directory to store suggested and expanded Kubernetes configuration files. Suggested files will be stored in \"<output>/suggested\" and expanded files will be stored in \"<output>/expanded\".")
	cmd.Flags().IntVarP(&options.exposePort, "expose", "x", 0, "Creates a Service object that connects to a deployed workload object using a selector that matches the label with key as 'app' and value of the image name's suffix specified by --image. The port provided will be used to expose the deployed workload object (i.e., port and targetPort will be set to the value provided in this flag).")
	cmd.Flags().BoolVarP(&options.verbose, "verbose", "V", false, "Prints underlying commands being called to stdout.")
	cmd.Flags().BoolVarP(&options.recursive, "recursive", "R", false, "Recursively search through the configuration directory for all yaml files.")

	return cmd
}

func prepare(_ *cobra.Command, options *options) error {
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
	if options.output == "" {
		return fmt.Errorf("value of -o|--output cannot be empty")
	}

	if options.exposePort < 0 {
		return fmt.Errorf("value of -x|--expose must be > 0")
	}
	if options.exposePort > 0 && options.image == "" {
		return fmt.Errorf("exposing a deployed workload object requires -i|--image to be set")
	}

	labelsMap, err := common.CreateMapFromEqualDelimitedStrings(options.labels)
	if err != nil {
		return err
	}
	annotationsMap, err := common.CreateMapFromEqualDelimitedStrings(options.annotations)
	if err != nil {
		return err
	}
	d, err := common.CreateDeployer(ctx, false /* useGcloud */, options.verbose)
	if err != nil {
		return err
	}

	if err := d.Prepare(ctx, im, options.appName, options.appVersion, options.filename, common.SuggestedOutputPath(options.output), common.ExpandedOutputPath(options.output), options.namespace, labelsMap, annotationsMap, options.exposePort, options.recursive); err != nil {
		return fmt.Errorf("failed to prepare deployment: %v", err)
	}

	return nil
}
