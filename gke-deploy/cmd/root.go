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
  - Modify Kubernetes config YAMLs:
    - Set the digest of images that match the [--image|-i] flag, if provided.
    - Add app.kubernetes.io/name=[--name|-a] label, if provided.
    - Add app.kubernetes.io/version=[--version|-v] label, if provided.

Apply Phase:
  - Apply Kubernetes config YAMLs to the target cluster with the provided namespace.
  - Wait for deployed resources to be ready before exiting.
`
	example = `  # Modify configs and deploy to GKE cluster.
  gke-deploy run -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o modified -n my-namespace -c my-cluster -l us-east1-b

  # Deploy to GKE cluster that kubectl is currently targeting.
  gke-deploy run -f configs

  # Prepare only.
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o modified -n my-namespace

  # Apply only.
  gke-deploy apply -f configs -c my-cluster -n my-namespace -c my-cluster -l us-east1-b

  # Execute prepare and apply, with an intermediary step in between (e.g., manually check modified YAMLs)
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o modified -n my-namespace
  cat modified/*
  gke-deploy apply -f modified -c my-cluster -n my-namespace -c my-cluster -l us-east1-b  # Pass modified directory to -f`
	version = "" // TODO(joonlim): Create plan for versioning.
)

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

func Execute() error {
	return NewCommand().Execute()
}
