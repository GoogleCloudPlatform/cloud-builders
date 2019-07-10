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
package deployer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/image"

	"github.com/google/go-containerregistry/pkg/name"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/cluster"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/gcp"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/resource"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

const (
	appNameLabelKey    = "app.kubernetes.io/name"
	appVersionLabelKey = "app.kubernetes.io/version"
	managedByLabelKey  = "app.kubernetes.io/managed-by"

	managedByLabelValue = "gcp-cloud-build-deploy"
)

// Deployer handles the deployment of an image to a cluster.
type Deployer struct {
	Clients *services.Clients
}

// Prepare handles preparing deployment.
func (d *Deployer) Prepare(ctx context.Context, images []name.Reference, appName, appVersion, config, output, namespace string, labels map[string]string) error {
	fmt.Printf("Preparing deployment.\n")

	objs, err := resource.ParseConfigs(ctx, config, d.Clients.OS)
	if err != nil {
		return fmt.Errorf("failed to parse configs %q: %v", config, err)
	}
	fmt.Printf("Configs to prepare: %v\n", objs)

	for _, im := range images {
		imageDigest, err := image.ResolveDigest(ctx, im, d.Clients.Remote)
		if err != nil {
			return fmt.Errorf("failed to get image digest: %v", err)
		}
		imageName := image.Name(im)
		if err != nil {
			return fmt.Errorf("failed to get image name: %v", err)
		}
		imageWithDigest := fmt.Sprintf("%s@%s", imageName, imageDigest)

		fmt.Printf("Got digest for image: %s --> %s\n", im, imageWithDigest)
		fmt.Printf("Updating resource containers that have image name %q\n", imageName)

		if err := resource.UpdateMatchingContainerImage(ctx, objs, imageName, imageWithDigest); err != nil {
			return fmt.Errorf("failed to update container of objects: %v", err)
		}
	}

	if namespace != "default" {
		ok, err := resource.HasObject(ctx, objs, "Namespace", namespace)
		if err != nil {
			return fmt.Errorf("failed to check if namespace %q exists: %v", namespace, err)
		}
		if !ok {
			fmt.Printf("Creating namespace resource %q\n", namespace)
			nsObj, err := resource.CreateNamespaceObject(ctx, namespace)
			if err != nil {
				return fmt.Errorf("failed to create namespace object: %v", err)
			}
			if err = resource.AddObject(ctx, objs, nsObj); err != nil {
				return fmt.Errorf("failed to add namespace object: %v", err)
			}
		}
	}

	fmt.Printf("Hydrating resources.\n")

	if err := resource.UpdateNamespace(ctx, objs, namespace); err != nil {
		return fmt.Errorf("failed to update namespace: %v", err)
	}

	for _, obj := range objs {
		if resource.ResourceKind(obj) != "Namespace" {
			if appName != "" {
				if err := resource.AddLabel(ctx, obj, appNameLabelKey, appName, false); err != nil {
					return fmt.Errorf("failed to add %s=%s label to object %v: %v", appNameLabelKey, appName, obj, err)
				}
			}
			if appVersion != "" {
				if err := resource.AddLabel(ctx, obj, appVersionLabelKey, appVersion, false); err != nil {
					return fmt.Errorf("failed to add %s=%s label to object %v: %v", appVersionLabelKey, appVersion, obj, err)
				}
			}
		}

		if err := resource.AddLabel(ctx, obj, managedByLabelKey, managedByLabelValue, true); err != nil {
			return fmt.Errorf("failed to add %s=%s label to object %v: %v", managedByLabelKey, managedByLabelValue, obj, err)
		}

		for k, v := range labels {
			if k == appNameLabelKey {
				return fmt.Errorf("%s label must be set using the --app|-a flag", appNameLabelKey)
			}
			if k == appVersionLabelKey {
				return fmt.Errorf("%s label must be set using the --version|-v flag", appVersionLabelKey)
			}
			if k == managedByLabelKey {
				return fmt.Errorf("%s label cannot be explicitly set", managedByLabelKey)
			}

			if err := resource.AddLabel(ctx, obj, k, v, true); err != nil {
				return fmt.Errorf("failed to add %s=%s custom label to object %v: %v", k, v, obj, err)
			}
		}
	}

	fmt.Printf("Saving hydrated resource configs to output: %q\n", output)
	if err := resource.SaveAsConfigs(ctx, objs, output, d.Clients.OS); err != nil {
		return fmt.Errorf("failed to save hydrated configs to output: %v", err)
	}

	fmt.Printf("Finished preparing deployment.\n\n")

	return nil
}

// Apply handles applying the deployment.
func (d *Deployer) Apply(ctx context.Context, clusterName, clusterLocation, clusterProject, config, namespace string, waitTimeout time.Duration) error {
	fmt.Printf("Applying deployment.\n")

	if (clusterName != "" && clusterLocation == "") || (clusterName == "" && clusterLocation != "") {
		return fmt.Errorf("clusterName and clusterLocation either must both be provided, or neither should be provided")
	}
	if clusterProject == "" {
		currentProject, err := gcp.GetProject(ctx, d.Clients.Gcloud)
		if err != nil {
			return fmt.Errorf("failed to get GCP project: %v", err)
		}
		clusterProject = currentProject
	}

	if clusterName != "" && clusterLocation != "" {
		fmt.Printf("Getting access to cluster %q in %q.\n", clusterName, clusterLocation)
		if err := cluster.AuthorizeAccess(ctx, clusterName, clusterLocation, clusterProject, d.Clients.Gcloud); err != nil {
			account, err2 := gcp.GetAccount(ctx, d.Clients.Gcloud)
			if err2 != nil {
				fmt.Printf("Failed to get GCP account. Swallowing error: %v\n", err)
			}
			if err2 == nil {
				// TODO(joonlim): Find a better way to figure out if accountType is "user", "serviceAccount", or "group".
				accountType := "user"
				if strings.Contains(account, "gserviceaccount.com") {
					accountType = "serviceAccount"
				}

				fmt.Printf("> You may need to grant permission to access to the cluster:\n\n")
				fmt.Printf("   gcloud projects add-iam-policy-binding %s --member=%s:%s --role=roles/container.developer\n\n", clusterProject, accountType, account)
			}
			return fmt.Errorf("failed to get access to cluster: %v", err)
		}
	}

	objs, err := resource.ParseConfigs(ctx, config, d.Clients.OS)
	if err != nil {
		return fmt.Errorf("failed to parse configs: %v", err)
	}
	fmt.Printf("Configs to apply: %v\n", objs)

	exists := make(map[string]bool)
	var dups []string
	for _, obj := range objs {
		key := fmt.Sprintf("%v", obj)
		ok := exists[key]
		if ok {
			dups = append(dups, key)
		}
		exists[key] = true
	}
	if len(dups) > 0 {
		fmt.Fprintf(os.Stderr, "\nWARNING: Deploying multiple resources share the same kind and name. Duplicate resources will be overridden:\n")
		for _, obj := range dups {
			fmt.Fprintf(os.Stderr, "%v\n", obj)
		}
		fmt.Fprintln(os.Stderr)
	}

	fmt.Printf("Applying resource configs to cluster.\n")

	// Apply all namespace objects first, if they exists
	for baseName, obj := range objs {
		if resource.ResourceKind(obj) == "Namespace" {
			nsFile := filepath.Join(config, baseName)
			if err := cluster.ApplyConfigs(ctx, nsFile, "", d.Clients.Kubectl); err != nil {
				return fmt.Errorf("failed to apply namespace config to cluster: %v", err)
			}
			// TODO(joonlim): Wait for deployed namespace to be ready before applying other objects
		}
	}

	if err := cluster.ApplyConfigs(ctx, config, namespace, d.Clients.Kubectl); err != nil {
		return fmt.Errorf("failed to apply configs to cluster: %v", err)
	}

	deployedObjs := resource.Objects{}
	timedOut := false

	fmt.Printf("Waiting for deployed objects to be ready with timeout of %v\n", waitTimeout)
	start := time.Now()
	end := start.Add(waitTimeout)
	periodicMsgInterval := 30 * time.Second
	nextPeriodicMsg := time.Now().Add(periodicMsgInterval)
	for len(objs) > 0 {
		for key, obj := range objs {
			kind := resource.ResourceKind(obj)
			name, err := resource.ResourceName(obj)
			if err != nil {
				return fmt.Errorf("failed to get name of resource: %v", err)
			}
			deployedObj, err := cluster.GetDeployedObject(ctx, kind, name, namespace, d.Clients.Kubectl)
			if err != nil {
				return fmt.Errorf("failed to get config of deployed object with kind %q and name %q: %v", kind, name, err)
			}
			deployedObjs[key] = deployedObj
			ok, err := resource.IsReady(ctx, deployedObj)
			if err != nil {
				return fmt.Errorf("failed to check if deployed object with kind %q and name %q is ready: %v", kind, name, err)
			}
			if ok {
				dur := time.Now().Sub(start).Round(time.Second / 10) // Round to nearest 0.1 seconds
				fmt.Printf("Deployed object with kind %q and name %q is ready after %v\n", kind, name, dur)
				delete(objs, key)
			}
		}
		if time.Now().After(end) {
			timedOut = true
			break
		}
		if time.Now().After(nextPeriodicMsg) {
			fmt.Printf("Still waiting on %d object(s) to be ready: %v\n", len(objs), objs)
			nextPeriodicMsg = nextPeriodicMsg.Add(periodicMsgInterval)
		}
		select {
		case <-time.After(5 * time.Second):
		}
	}

	fmt.Printf("Finished applying deployment.\n\n")

	summary, err := resource.DeploySummary(ctx, deployedObjs)
	if err != nil {
		return fmt.Errorf("failed to get summary of deployed resources: %v", err)
	}

	fmt.Printf("################################################################################\n")
	fmt.Printf("> Deployed Resources\n\n")
	fmt.Printf("%s\n", summary)

	fmt.Printf("################################################################################\n")

	links, err := d.gkeLinks(clusterProject)
	if err != nil {
		return fmt.Errorf("failed to get GKE links: %v", err)
	}

	fmt.Printf("> GKE\n\n")
	fmt.Printf("%s\n", links)

	if timedOut {
		return fmt.Errorf("timed out after %v while waiting for deployed objects to be ready", waitTimeout)
	}

	return nil
}

func (d *Deployer) gkeLinks(clusterProject string) (string, error) {
	padding := 4
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, padding, ' ', 0)

	if _, err := fmt.Fprintf(w, "Workloads:\thttps://console.cloud.google.com/kubernetes/workload?project=%s\n", clusterProject); err != nil {
		return "", fmt.Errorf("failed to write to writer: %v", err)
	}
	if _, err := fmt.Fprintf(w, "Services & Ingress:\thttps://console.cloud.google.com/kubernetes/discovery?project=%s\n", clusterProject); err != nil {
		return "", fmt.Errorf("failed to write to writer: %v", err)
	}
	if _, err := fmt.Fprintf(w, "Applications:\thttps://console.cloud.google.com/kubernetes/applications?project=%s\n", clusterProject); err != nil {
		return "", fmt.Errorf("failed to write to writer: %v", err)
	}
	if _, err := fmt.Fprintf(w, "Configuration:\thttps://console.cloud.google.com/kubernetes/config?project=%s\n", clusterProject); err != nil {
		return "", fmt.Errorf("failed to write to writer: %v", err)
	}
	if _, err := fmt.Fprintf(w, "Storage:\thttps://console.cloud.google.com/kubernetes/storage?project=%s\n", clusterProject); err != nil {
		return "", fmt.Errorf("failed to write to writer: %v", err)
	}

	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush writer: %v", err)
	}

	return buf.String(), nil
}
