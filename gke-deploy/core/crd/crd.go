// Package crd contains logic related to CRDs.
package crd

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

const (
	applicationCRDName       = "customresourcedefinition.apiextensions.k8s.io/applications.app.k8s.io"
	applicationCRDInstallURI = "https://raw.githubusercontent.com/kubernetes-sigs/application/master/config/crd/bases/app.k8s.io_applications.yaml"
)

// EnsureInstallApplicationCRD ensures the installation of the Application CRD in the current context's cluster.
func EnsureInstallApplicationCRD(ctx context.Context, ks services.KubectlService) error {
	installed, err := crdIsInstalled(ctx, applicationCRDName, ks)
	if err != nil {
		return err
	}
	if installed {
		return nil
	}
	if err := ks.Apply(ctx, applicationCRDInstallURI, ""); err != nil {
		return fmt.Errorf("failed to apply Application CRD: %v", err)
	}
	return nil
}

// crdIsInstalled returns true if a CRD <crd> is installed in the current context's cluster,
// else false.
func crdIsInstalled(ctx context.Context, crd string, ks services.KubectlService) (bool, error) {
	objYaml, err := ks.Get(ctx, crd, "", "", "yaml", true)
	if err != nil {
		return false, fmt.Errorf("failed to get config of CRD %q: %v", applicationCRDName, err)
	}
	if objYaml == "" {
		return false, nil
	}
	return true, nil
}
