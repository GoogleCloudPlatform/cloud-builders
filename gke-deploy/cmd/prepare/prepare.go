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
package prepare

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/common"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/image"
)

const (
	short = "Execute prepare phase and skip apply phase"
	long  = `Prepare to deploy to GKE by generating modified Kubernetes resource configs. Skip apply.

- Modify Kubernetes config YAMLs to:
  - Set the digest of images that match the [--image|-i] flag, if provided.
  - Add app.kubernetes.io/name=[--app|-a] label, if provided.
  - Add app.kubernetes.io/version=[--version|-v] label, if provided.
`
	example = `  # Prepare only.
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o modified -n my-namespace

  # Execute prepare and apply, with an intermediary step in between (e.g., manually check modified YAMLs)
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o modified -n my-namespace
  cat modified/*
  gke-deploy apply -f modified -c my-cluster -n my-namespace -c my-cluster -l us-east1-b  # Pass modified directory to -f`
)

type options struct {
	appName    string
	appVersion string
	filename   string
	images     []string
	labels     []string
	namespace  string
	output     string
	verbose    bool
}

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
	cmd.Flags().StringVarP(&options.filename, "filename", "f", "", "Config file or directory of config files to use to create the Kubernetes resources (file or files in directory must end with \".yml\" or \".yaml\").")
	cmd.Flags().StringSliceVarP(&options.images, "image", "i", nil, "Image(s) to be deployed. Images can be set comma-delimited or as separate flags.")
	cmd.Flags().StringSliceVarP(&options.labels, "label", "L", nil, "Label(s) to add to Kubernetes resources (k1=v1). Labels can be set comma-delimited or as separate flags. If two or more labels with the same key are listed, the last one is used.")
	cmd.Flags().StringVarP(&options.namespace, "namespace", "n", "default", "Name of GKE cluster to deploy to.")
	cmd.Flags().StringVarP(&options.output, "output", "o", "./output", "Target directory to store modified Kubernetes resource configs.")
	cmd.Flags().BoolVarP(&options.verbose, "verbose", "V", false, "Prints underlying commands being called to stdout.")

	return cmd
}

func prepare(cmd *cobra.Command, options *options) error {
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

	return nil
}
