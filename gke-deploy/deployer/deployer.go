// Package deployer contains logic related to deploying to a GKE cluster.
package deployer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/crd"
	"github.com/google/go-containerregistry/pkg/name"
	applicationsv1beta1 "github.com/kubernetes-sigs/application/pkg/apis/app/v1beta1"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/cluster"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/gcp"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/gcs"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/image"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/core/resource"
	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/services"
)

const (
	appNameLabelKey    = "app.kubernetes.io/name"
	appVersionLabelKey = "app.kubernetes.io/version"
	managedByLabelKey  = "app.kubernetes.io/managed-by"

	managedByLabelValue = "gcp-cloud-build-deploy"
	// Name pattern used to create temporary staging folders for the files to be downloaded from/uploaded to GCS.
	k8sConfigStagingDir = "gke_deploy_temp_"
	expendedFileName    = "expanded-resources.yaml"
	suggestedFileName   = "suggested-resources.yaml"
)

// Deployer handles the deployment of an image to a cluster.
type Deployer struct {
	Clients      *services.Clients
	UseGcloud    bool
	ServerDryRun bool
}

// Prepare handles preparing deployment.
func (d *Deployer) Prepare(ctx context.Context, im name.Reference, appName, appVersion, config, suggestedOutput, expandedOutput, namespace string, labels, annotations map[string]string, exposePort int, recursive, createApplicationCR bool, applicationLinks []applicationsv1beta1.Link) error {
	fmt.Printf("Preparing deployment.\n")

	var objs resource.Objects
	ss := &gcs.GCS{
		GcsService: d.Clients.GCS,
	}
	if config != "" {

		if strings.HasPrefix(config, "gs://") {
			tmpDir, err := d.Clients.OS.TempDir(ctx, "", k8sConfigStagingDir)
			if err != nil {
				return fmt.Errorf("failed to create tmp directory: %v", err)
			}
			defer d.Clients.OS.RemoveAll(ctx, tmpDir)
			err = ss.Download(ctx, config, tmpDir, recursive)
			if err != nil {
				return fmt.Errorf("failed to download configuration files from GCS %q: %v", config, err)
			}
			config = tmpDir
		}

		parsed, err := resource.ParseConfigs(ctx, config, d.Clients.OS, recursive)
		if err != nil {
			return fmt.Errorf("failed to parse configuration files %q: %v", config, err)
		}
		if len(parsed) == 0 {
			return fmt.Errorf("no objects found")
		}
		objs = parsed
		fmt.Printf("Configuration files to be used: %v\n", objs)
	} else {
		objs = resource.Objects{}
		fmt.Println("Starting with no configuration files")
	}

	if im != nil {
		// e.g., Resolve "gcr.io/my-project/my-app:1.0.0" to name suffix "my-app".
		imageNameSplit := strings.Split(image.Name(im), "/")
		imageNameSuffix := imageNameSplit[len(imageNameSplit)-1]
		imageName := image.Name(im)

		if config == "" {
			fmt.Printf("Creating suggested Deployment configuration file %q\n", imageNameSuffix)
			dObj, err := resource.CreateDeploymentObject(ctx, imageNameSuffix, imageNameSuffix, imageName)
			if err != nil {
				return fmt.Errorf("failed to create Deployment object: %v", err)
			}

			objs = append(objs, dObj)

			hpaName := fmt.Sprintf("%s-hpa", imageNameSuffix)
			fmt.Printf("Creating suggested HorizontalPodAutoscaler configuration file %q\n", hpaName)
			hpaObj, err := resource.CreateHorizontalPodAutoscalerObject(ctx, hpaName, imageNameSuffix)
			if err != nil {
				return fmt.Errorf("failed to create HorizontalPodAutoscaler object: %v", err)
			}
			objs = append(objs, hpaObj)
		}

		// Remove tag/digest from image references.
		if err := resource.UpdateMatchingContainerImage(ctx, objs, imageName, imageName); err != nil {
			return fmt.Errorf("failed to update container of objects: %v", err)
		}
	}

	if appName != "" {
		if exposePort > 0 {
			service := fmt.Sprintf("%s-service", appName)
			ok, err := resource.HasObject(ctx, objs, "Service", service)
			if err != nil {
				return fmt.Errorf("failed to check if Service %q exists: %v", service, err)
			}
			if !ok {
				fmt.Printf("Creating suggested Service configuration file %q\n", service)
				svcObj, err := resource.CreateServiceObject(ctx, service, appNameLabelKey, appName, exposePort)
				if err != nil {
					return fmt.Errorf("failed to create Service object: %v", err)
				}
				objs = append(objs, svcObj)
			} else {
				fmt.Fprintf(os.Stderr, "\nWARNING: Service %q already exists in provided configuration files. Not generating new Service.\n\n", service)
			}
		}

		if createApplicationCR {
			ok, err := resource.HasObject(ctx, objs, "Application", appName)
			if err != nil {
				return fmt.Errorf("failed to check if Application %q exists: %v", appName, err)
			}
			if !ok {
				fmt.Printf("Creating suggested Application configuration file %q\n", appName)
				appObj, err := resource.CreateApplicationObject(appName, appNameLabelKey, appName, appName, appVersion, objs)
				if err != nil {
					return fmt.Errorf("failed to create Application object: %v", err)
				}
				objs = append(objs, appObj)
			} else {
				fmt.Fprintf(os.Stderr, "\nWARNING: Application %q already exists in provided configuration files. Not generating new Application.\n\n", appName)
			}
		}
	}

	if namespace != "" && namespace != "default" {
		ok, err := resource.HasObject(ctx, objs, "Namespace", namespace)
		if err != nil {
			return fmt.Errorf("failed to check if Namespace %q exists: %v", namespace, err)
		}
		if !ok {
			fmt.Printf("Creating suggested Namespace configuration file %q\n", namespace)
			nsObj, err := resource.CreateNamespaceObject(ctx, namespace)
			if err != nil {
				return fmt.Errorf("failed to create Namespace object: %v", err)
			}
			objs = append(objs, nsObj)
		}
	}

	for _, obj := range objs {
		if resource.ObjectKind(obj) != "Namespace" {
			if appName != "" {
				if err := resource.AddLabel(ctx, obj, appNameLabelKey, appName, false); err != nil {
					return fmt.Errorf("failed to add %s=%s label to object %v: %v", appNameLabelKey, appName, obj, err)
				}
			}
		}
	}

	var toGcs bool
	var gcsPath string
	if len(objs) > 0 {
		fmt.Printf("Saving suggested configuration files to %q\n", suggestedOutput)
		var lineComments map[string]string
		if im != nil {
			lineComments = map[string]string{
				fmt.Sprintf("image: %s", image.Name(im)): "Will be set to actual image before deployment",
			}
		}

		if strings.HasPrefix(suggestedOutput, "gs://") {
			tmpDir, err := d.Clients.OS.TempDir(ctx, "", k8sConfigStagingDir)
			if err != nil {
				return fmt.Errorf("failed to create tmp directory: %v", err)
			}
			defer d.Clients.OS.RemoveAll(ctx, tmpDir)
			gcsPath = strings.Join([]string{suggestedOutput, suggestedFileName}, "/")
			suggestedOutput = tmpDir
			toGcs = true
		}

		fileName, err := resource.SaveAsConfigs(ctx, objs, suggestedOutput, lineComments, d.Clients.OS)
		if err != nil {
			return fmt.Errorf("failed to save suggested configuration files to %q: %v", suggestedOutput, err)
		}

		if toGcs {
			err := ss.Upload(ctx, fileName, gcsPath)
			if err != nil {
				return fmt.Errorf("failed to upload configuration files from GCS %q: %v", config, err)
			}
		}
	}

	fmt.Printf("\nExpanding configuration files.\n")

	if im != nil {
		imageName := image.Name(im)
		imageDigest, err := image.ResolveDigest(ctx, im, d.Clients.Remote)
		if err != nil {
			return fmt.Errorf("failed to get image digest: %v", err)
		}
		imageWithDigest := fmt.Sprintf("%s@%s", image.Name(im), imageDigest)
		fmt.Printf("Got digest for image: %s --> %s\n", im, imageWithDigest)

		fmt.Printf("Updating containers in configuration files that have image name %q to use image with digest %q\n", imageName, imageWithDigest)
		if err := resource.UpdateMatchingContainerImage(ctx, objs, imageName, imageWithDigest); err != nil {
			return fmt.Errorf("failed to update container of objects: %v", err)
		}
	}

	if namespace != "" {
		if err := resource.UpdateNamespace(ctx, objs, namespace); err != nil {
			return fmt.Errorf("failed to update namespace of objects: %v", err)
		}
	} else {
		if err := resource.AddNamespaceIfMissing(objs, "default"); err != nil {
			return fmt.Errorf("failed to update namespace of objects with no namespace to default: %v", err)
		}
	}

	for _, obj := range objs {
		if resource.ObjectKind(obj) != "Namespace" {
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

		for k, v := range annotations {
			if err := resource.AddAnnotation(obj, k, v); err != nil {
				return fmt.Errorf("failed to add %s=%s custom annotation to object %v: %v", k, v, obj, err)
			}
		}

		if resource.ObjectKind(obj) == "Application" {
			if err := resource.SetApplicationLinks(obj, applicationLinks); err != nil {
				return fmt.Errorf("failed to add links to Application: %v", err)
			}
		}
	}

	fmt.Printf("Saving expanded configuration files to %q\n", expandedOutput)

	if strings.HasPrefix(expandedOutput, "gs://") {
		tmpDir, err := d.Clients.OS.TempDir(ctx, "", k8sConfigStagingDir)
		if err != nil {
			return fmt.Errorf("failed to create tmp directory: %v", err)
		}
		defer d.Clients.OS.RemoveAll(ctx, tmpDir)
		gcsPath = strings.Join([]string{expandedOutput, expendedFileName}, "/")
		expandedOutput = tmpDir
		toGcs = true
	}

	fileName, err := resource.SaveAsConfigs(ctx, objs, expandedOutput, nil, d.Clients.OS)
	if err != nil {
		return fmt.Errorf("failed to save expanded configuration files to %q: %v", expandedOutput, err)
	}

	if toGcs {
		err := ss.Upload(ctx, fileName, gcsPath)
		if err != nil {
			return fmt.Errorf("failed to upload configuration files from GCS %q: %v", config, err)
		}
	}

	fmt.Printf("Finished preparing deployment.\n\n")

	return nil
}

// Apply handles applying the deployment.
func (d *Deployer) Apply(ctx context.Context, clusterName, clusterLocation, clusterProject, config, namespace string, waitTimeout time.Duration, recursive bool) error {
	if d.ServerDryRun {
		fmt.Printf("Applying deployment in server dry run mode.\n")
	} else {
		fmt.Printf("Applying deployment.\n")
	}

	if (clusterName != "" && clusterLocation == "") || (clusterName == "" && clusterLocation != "") {
		return fmt.Errorf("clusterName and clusterLocation either must both be provided, or neither should be provided")
	}
	if clusterProject == "" && d.UseGcloud {
		currentProject, err := gcp.GetProject(ctx, d.Clients.Gcloud)
		if err != nil {
			return fmt.Errorf("failed to get GCP project: %v", err)
		}
		clusterProject = currentProject
	}

	if clusterName != "" && clusterLocation != "" && d.UseGcloud {
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

	if strings.HasPrefix(config, "gs://") {

		tmpDir, err := d.Clients.OS.TempDir(ctx, "", k8sConfigStagingDir)
		if err != nil {
			return fmt.Errorf("failed to create tmp directory: %v", err)
		}
		defer d.Clients.OS.RemoveAll(ctx, tmpDir)
		ss := &gcs.GCS{
			GcsService: d.Clients.GCS,
		}
		err = ss.Download(ctx, config, tmpDir, recursive)
		if err != nil {
			return fmt.Errorf("failed to download configuration files from GCS %q: %v", config, err)
		}
		config = tmpDir
	}

	objs, err := resource.ParseConfigs(ctx, config, d.Clients.OS, recursive)
	if err != nil {
		return fmt.Errorf("failed to parse configuration files: %v", err)
	}
	if len(objs) == 0 {
		return fmt.Errorf("no objects found")
	}
	fmt.Printf("Configuration files to be used: %v\n", objs)

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
		fmt.Fprintf(os.Stderr, "\nWARNING: Deploying multiple objects share the same kind and name. Duplicate objects will be overridden:\n")
		for _, obj := range dups {
			fmt.Fprintf(os.Stderr, "%v\n", obj)
		}
		fmt.Fprintln(os.Stderr)
	}

	fmt.Printf("Applying configuration files to cluster.\n")

	// Apply all namespace objects first, if they exist
	filteredObjs := make(resource.Objects, 0, len(objs))
	for _, obj := range objs {
		if resource.ObjectKind(obj) == "Namespace" {
			nsName, err := resource.ObjectName(obj)
			if err != nil {
				return fmt.Errorf("failed to get name of object: %v", err)
			}
			exists, err := cluster.DeployedObjectExists(ctx, "Namespace", nsName, "", d.Clients.Kubectl)
			if err != nil {
				return fmt.Errorf("failed to check if deployed object with kind \"Namespace\" and name %q exists: %v", nsName, err)
			}
			if !exists {
				fmt.Fprintf(os.Stderr, "\nWARNING: It is recommended that namespaces be created by an administrator. Creating namespace %q because it does not exist.\n\n", nsName)
				objString, err := resource.EncodeToYAMLString(obj)
				if err != nil {
					return fmt.Errorf("failed to encode obj to string")
				}
				if err := cluster.ApplyConfigFromString(ctx, objString, "", d.Clients.Kubectl); err != nil {
					return fmt.Errorf("failed to apply Namespace configuration file with name %q to cluster: %v", nsName, err)
				}
			}
		} else {
			// Delete namespace from list of objects to be deployed because it has already been deployed we do not want it to show up in the deployment summary.
			filteredObjs = append(filteredObjs, obj)
		}
	}

	objs = filteredObjs

	// Apply each config file individually vs applying the directory to avoid applying namespaces.
	// Namespace objects are removed from objs at this point.
	ensuredInstallApplicationCRD := false // Only need to do this once, in the case where the user provides more than one Application CR
	for _, obj := range objs {
		objName, err := resource.ObjectName(obj)
		if err != nil {
			return fmt.Errorf("failed to get name of object: %v", err)
		}

		if !ensuredInstallApplicationCRD && resource.ObjectKind(obj) == "Application" {
			if err := crd.EnsureInstallApplicationCRD(ctx, d.Clients.Kubectl); err != nil {
				return fmt.Errorf("failed to ensure installation of Application CRD on target cluster: %v", err)
			}
			ensuredInstallApplicationCRD = true
		}

		objString, err := resource.EncodeToYAMLString(obj)
		if err != nil {
			return fmt.Errorf("failed to encode obj to string")
		}
		// If namespace == "", uses the namespace defined in each config.
		if err := cluster.ApplyConfigFromString(ctx, objString, namespace, d.Clients.Kubectl); err != nil {
			return fmt.Errorf("failed to apply %s configuration file with name %q to cluster: %v", resource.ObjectKind(obj), objName, err)
		}
	}

	deployedObjs := map[string]map[string]resource.Object{}
	summaryObjs := make(resource.Objects, 0, len(objs))
	timedOut := false

	if d.ServerDryRun {
		fmt.Printf("Server-side dry run deployment succeeded.\n\n")
		return nil
	}

	fmt.Printf("\nWaiting for deployed objects to be ready with timeout of %v\n", waitTimeout)
	start := time.Now()
	end := start.Add(waitTimeout)
	periodicMsgInterval := 30 * time.Second
	nextPeriodicMsg := time.Now().Add(periodicMsgInterval)
	ticker := time.NewTicker(5 * time.Second)
	for len(objs) > 0 {

		filteredObjs := make(resource.Objects, 0, len(objs))

		for _, obj := range objs {
			kind := resource.ObjectKind(obj)
			name, err := resource.ObjectName(obj)
			if err != nil {
				return fmt.Errorf("failed to get name of object: %v", err)
			}
			objNamespace := ""
			if namespace == "" {
				ns, err := resource.ObjectNamespace(obj)
				if err != nil {
					return fmt.Errorf("failed to get namespace of object: %v", err)
				}
				objNamespace = ns
			} else {
				objNamespace = namespace
			}
			deployedObj, err := cluster.GetDeployedObject(ctx, kind, name, objNamespace, d.Clients.Kubectl)
			if err != nil {
				return fmt.Errorf("failed to get configuration of deployed object with kind %q and name %q: %v", kind, name, err)
			}
			if deployedObjs[kind] == nil {
				deployedObjs[kind] = map[string]resource.Object{}
			}
			deployedObjs[kind][name] = *deployedObj
			ok, err := resource.IsReady(ctx, deployedObj)
			if err != nil {
				return fmt.Errorf("failed to check if deployed object with kind %q and name %q is ready: %v", kind, name, err)
			}
			if ok {
				dur := time.Now().Sub(start).Round(time.Second / 10) // Round to nearest 0.1 seconds
				fmt.Printf("Deployed object with kind %q and name %q is ready after %v\n", kind, name, dur)
			} else {
				filteredObjs = append(filteredObjs, obj)
			}
		}

		objs = filteredObjs

		if len(objs) == 0 {
			// Break out here to avoid waiting for ticker.
			break
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
		case <-ticker.C:
		}
	}

	fmt.Printf("Finished applying deployment.\n\n")

	for _, nameMap := range deployedObjs {
		for k, _ := range nameMap {
			o := nameMap[k]
			summaryObjs = append(summaryObjs, &o)
		}
	}
	summary, err := resource.DeploySummary(ctx, summaryObjs)
	if err != nil {
		return fmt.Errorf("failed to get summary of deployed objects: %v", err)
	}

	fmt.Printf("################################################################################\n")
	fmt.Printf("> Deployed Objects\n\n")
	fmt.Printf("%s\n", summary)

	fmt.Printf("################################################################################\n")

	if clusterProject != "" {
		links, err := d.gkeLinks(clusterProject)
		if err != nil {
			return fmt.Errorf("failed to get GKE links: %v", err)
		}

		fmt.Printf("> GKE\n\n")
		fmt.Printf("%s\n", links)
	}

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
	if _, err := fmt.Fprintf(w, "Applications:\thttps://console.cloud.google.com/kubernetes/application?project=%s\n", clusterProject); err != nil {
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
