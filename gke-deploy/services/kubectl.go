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
package services

import (
	"context"
	"fmt"
	"os/exec"
)

// Kubectl implements the KubectlService interface.
// The service account that is calling this must have permission to access the cluster.
// e.g., to run on GCB: gcloud projects add-iam-policy-binding <project-id> --member=serviceAccount:<project-number>@cloudbuild.gserviceaccount.com --role=roles/container.admin
type Kubectl struct {
	printCommands bool
}

// NewKubectl returns a new Kubectl object.
func NewKubectl(ctx context.Context, printCommands bool) (*Kubectl, error) {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return nil, err
	}
	return &Kubectl{
		printCommands,
	}, nil
}

// ApplyFromString calls `kubectl apply -f - -n <namespace> < ${configString}`.
func (k *Kubectl) ApplyFromString(configString, namespace string) error {
	args := []string{"apply", "-f", "-"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if _, err := runCommandWithStdinRedirection(k.printCommands, "kubectl", configString, args...); err != nil {
		return fmt.Errorf("command to apply kubernetes config from string to cluster failed: %v", err)
	}
	return nil
}

// Get calls `kubectl get <kind> <name> -n <namespace> --output=<format>`.
func (k *Kubectl) Get(ctx context.Context, kind, name, namespace, format string, ignoreNotFound bool) (string, error) {
	args := []string{"get", kind, name}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if format != "" {
		args = append(args, fmt.Sprintf("--output=%s", format))
	}
	if ignoreNotFound {
		args = append(args, "--ignore-not-found=true")
	}
	out, err := runCommand(k.printCommands, "kubectl", args...)
	if err != nil {
		return "", fmt.Errorf("command to get kubernetes config: %v", err)
	}
	return out, nil
}
