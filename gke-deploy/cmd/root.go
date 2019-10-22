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
// Package cmd contains the logic for `gke-deploy` top-level command.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/apply"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/prepare"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/run"
)

const (
	short = "Deploy to GKE"
	long  = `Deploy to GKE in two phases, which will do the following:

Prepare Phase:
  - Expand Kubernetes configuration files:
    - Set the digest of images that match the [--image|-i] flag, if provided.
    - Add app.kubernetes.io/name=[--app|-a] label, if provided.
    - Add app.kubernetes.io/version=[--version|-v] label, if provided.

Apply Phase:
  - Apply Kubernetes configuration files to the target cluster with the provided namespace.
  - Wait for deployed Kubernetes configuration files to be ready before exiting.
`
	example = `  # Expand Kubernetes configuration files and deploy to GKE cluster.
  gke-deploy run -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o expanded -n my-namespace -c my-cluster -l us-east1-b

  # Deploy to GKE cluster that kubectl is currently targeting.
  gke-deploy run -f configs

  # Deploy to GKE cluster that kubectl is currently targeting without supplying any Kubernetes configuration files. Have gke-deploy generate suggested Kubernetes configuration files for your application using an image, app name, and service port.
  gke-deploy run -i nginx -a nginx -x 80

  # Prepare only.
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o expanded -n my-namespace

  # Apply only.
  gke-deploy apply -f configs -c my-cluster -n my-namespace -c my-cluster -l us-east1-b

  # Execute prepare and apply, with an intermediary step in between (e.g., manually check expanded YAMLs)
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o expanded -n my-namespace
  cat expanded/*
  gke-deploy apply -f expanded -c my-cluster -n my-namespace -c my-cluster -l us-east1-b  # Pass expanded directory to -f

  # Pipe output from another templating engine to gke-deploy.
  kustomize build overlays/staging | gke-deploy run -f - -a my-app -c my-cluster -l us-east1-b
  helm template charts/prometheus | gke-deploy apply -f - -c my-cluster -l us-east1-b  # No need to run Tiller in cluster`
	version = "" // TODO(joonlim): Create plan for versioning.
)

// NewCommand creates the `gke-deploy` top-level command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gke-deploy",
		Short:   short,
		Long:    long,
		Example: example,
		Version: version,
	}

	cmd.AddCommand(apply.NewApplyCommand())
	cmd.AddCommand(prepare.NewPrepareCommand())
	cmd.AddCommand(run.NewRunCommand())

	return cmd
}

// Execute executes the `gke-deploy` top-level command.
func Execute() error {
	return NewCommand().Execute()
}
