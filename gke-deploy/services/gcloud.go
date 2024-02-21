package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	googleScopes = []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email"}
	kubeArgs = []string{"--use_application_default_credentials"}
)

const (
	gkeContextFormat      = "gke_%s_%s_%s"
	gkeResourceNameFormat = "projects/%s/locations/%s/clusters/%s"
	kubeApiVersion        = "client.authentication.k8s.io/v1beta1"
	kubeCommand           = "gke-gcloud-auth-plugin"
	kubeInstallHint       = `Install gke-gcloud-auth-plugin for use with kubectl by following
https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke`
	kubeProvideClusterInfo = true
)

// Gcloud implements the GcloudService interface.
type Gcloud struct {
	printCommands bool
}

// NewGcloud returns a new Gcloud object.
func NewGcloud(ctx context.Context, printCommands bool) (*Gcloud, error) {
	if _, err := exec.LookPath("gcloud"); err != nil {
		return nil, err
	}
	return &Gcloud{
		printCommands,
	}, nil
}

// NewGcloudGoClient returns a new Gcloud object.
func NewGcloudGoClient(ctx context.Context, printCommands bool) *Gcloud {
	return &Gcloud{
		printCommands,
	}
}

func (g *Gcloud) ContainerClustersGetCredentials(ctx context.Context, clusterName, clusterLocation, clusterProject string) error {
	if _, err := runCommand(ctx, g.printCommands, "gcloud", "container", "clusters", "get-credentials", clusterName, fmt.Sprintf("--zone=%s", clusterLocation), fmt.Sprintf("--project=%s", clusterProject), "--quiet"); err != nil {
		return fmt.Errorf("command to get cluster credentials failed: %v", err)
	}
	return nil
}

// ContainerClustersGetCredentialsGoClient uses the go client libraries to get cluster credentials and generate the kubeconfig file for kubectl.
// The kubectl authenticates using the accessToken instead of the google-gke-auth-plugin (which depends on gcloud).
func (g *Gcloud) ContainerClustersGetCredentialsGoClient(ctx context.Context, clusterName, clusterLocation, clusterProject string) error {

	dirname, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get the user's home directory: %v", err)
	}
	path := filepath.Join(dirname, ".kube")
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to make directory %s: %v", path, err)
	}
	kubeConfigFile := filepath.Join(path, "config")
	if err := getK8sClusterConfigs(ctx, clusterProject, clusterLocation, clusterName, kubeConfigFile); err != nil {
		return fmt.Errorf("failed to initialize gke cluster config: %v", err)
	}
	return nil
}

// ConfigGetValue calls `gcloud config get-value <property>` and returns stdout.
func (g *Gcloud) ConfigGetValue(ctx context.Context, property string) (string, error) {
	out, err := runCommand(ctx, g.printCommands, "gcloud", "config", "get-value", property, "--quiet")
	if err != nil {
		return "", fmt.Errorf("command to get property value failed: %v", err)
	}
	return strings.TrimSpace(out), nil
}

// getK8sClusterConfigs uses the go client libraries to authenticate against the cluster and generate the kubeconfig file
// instead of the gcloud CLI.
func getK8sClusterConfigs(ctx context.Context, clusterProject, clusterLocation, clusterName, kubeConfigFile string) error {
	fullClusterName := fmt.Sprintf(gkeResourceNameFormat, clusterProject, clusterLocation, clusterName)
	fmt.Printf("Full Cluster: %s\n", fullClusterName)
	ts, err := google.DefaultTokenSource(ctx, googleScopes...)
	if err != nil {
		return fmt.Errorf("failed to create google token source: %v", err)
	}
	options := []option.ClientOption{option.WithTokenSource(ts)}
	userAgent := os.Getenv("GOOGLE_APIS_USER_AGENT")
	if userAgent != "" {
		options = append(options, option.WithUserAgent(userAgent))
	}
	svc, err := container.NewService(ctx, options...)
	if err != nil {
		return fmt.Errorf("failed to initialize gke client: %v", err)
	}
	cluster, err := svc.Projects.Locations.Clusters.Get(fullClusterName).Do()
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}
	fmt.Printf("Cluster's Endpoint is: %s\n", cluster.Endpoint)
	name := fmt.Sprintf(gkeContextFormat, clusterProject, cluster.Zone, cluster.Name)
	kubeConfig := &api.Config{
		APIVersion:     "v1",
		Kind:           "Config",
		Clusters:       map[string]*api.Cluster{},
		AuthInfos:      map[string]*api.AuthInfo{},
		Contexts:       map[string]*api.Context{},
		CurrentContext: name,
	}
	_, err = os.Stat(kubeConfigFile)
	if os.IsNotExist(err) {
	} else {
		conf, err := clientcmd.LoadFromFile(kubeConfigFile)
		if err != nil {
			return fmt.Errorf("failed to load kubeConfig from file %s : %v", kubeConfigFile, err)
		}
		kubeConfig = conf
	}
	cert, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return fmt.Errorf("failed to decode the cluster certificate: %v", err)
	}
	kubeConfig.Clusters[name] = &api.Cluster{
		CertificateAuthorityData: cert,
		Server:                   "https://" + cluster.Endpoint,
	}
	kubeConfig.Contexts[name] = &api.Context{
		Cluster:  name,
		AuthInfo: name,
	}
	kubeConfig.AuthInfos[name] = &api.AuthInfo{
		Exec: &api.ExecConfig{
			APIVersion:         kubeApiVersion,
			Command:            kubeCommand,
			Args:               kubeArgs,
			InstallHint:        kubeInstallHint,
			ProvideClusterInfo: kubeProvideClusterInfo,
		},
	}
	if err := clientcmd.WriteToFile(*kubeConfig, kubeConfigFile); err != nil {
		return fmt.Errorf("failed to write kubeConfig to file %s: %v", kubeConfigFile, err)
	}
	return nil
}
