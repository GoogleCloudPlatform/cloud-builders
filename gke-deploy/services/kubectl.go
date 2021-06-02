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
	serverDryRun  bool
}

// NewKubectl returns a new Kubectl object.
func NewKubectl(ctx context.Context, printCommands bool, serverDryRun bool) (*Kubectl, error) {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return nil, err
	}
	return &Kubectl{
		printCommands,
		serverDryRun,
	}, nil
}

// Apply calls `kubectl apply -f <filename> n <namespace>`.
func (k *Kubectl) Apply(ctx context.Context, filename, namespace string) error {
	args := []string{"apply", "-f", filename}
	if k.serverDryRun {
		args = append(args, "--dry-run=server")
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if _, err := runCommand(ctx, k.printCommands, "kubectl", args...); err != nil {
		return fmt.Errorf("command to apply kubernetes config(s) to cluster failed: %v", err)
	}
	return nil
}

// ApplyFromString calls `kubectl apply -f - -n <namespace> < ${configString}`.
func (k *Kubectl) ApplyFromString(ctx context.Context, configString, namespace string) error {
	args := []string{"apply", "-f", "-"}
	if k.serverDryRun {
		args = append(args, "--dry-run=server")
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if _, err := runCommandWithStdinRedirection(ctx, k.printCommands, "kubectl", configString, args...); err != nil {
		return fmt.Errorf("command to apply kubernetes config from string to cluster failed: %v", err)
	}
	return nil
}

// Get calls `kubectl get <kind> <name> -n <namespace> --output=<format>`.
func (k *Kubectl) Get(ctx context.Context, kind, name, namespace, format string, ignoreNotFound bool) (string, error) {
	args := []string{"get", kind}
	if name != "" {
		args = append(args, name)
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if format != "" {
		args = append(args, fmt.Sprintf("--output=%s", format))
	}
	if ignoreNotFound {
		args = append(args, "--ignore-not-found=true")
	}
	out, err := runCommand(ctx, k.printCommands, "kubectl", args...)
	if err != nil {
		return "", fmt.Errorf("command to get kubernetes config: %v", err)
	}
	return out, nil
}
