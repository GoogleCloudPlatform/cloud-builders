# gke-deploy apply [flags]

This command will execute only the [Apply](../README.md#apply-step) phase.

## Parameters

Flag | Default Value | Description
--- | --- | --- |
`--filename`&#124;`-f` | | Config file or directory of config files to use to create the Kubernetes resources (file or files in directory must end with \".yml\" or \".yaml\"). If this is omitted, base YAMLs will be generated for the user.
`--namespace`&#124;`-n` | `default` | Kubernetes namespace in which to deploy the app.
`--cluster`&#124;`-c` | | Name of GKE cluster to deploy to.
`--location`&#124;`-l` | | Region/zone of GKE cluster to deploy to.
`--project`&#124;`-p` | | Project of GKE cluster to deploy  to. If this field is not provided the current set GCP project is used.
`--timeout`&#124;`-t` | 5m | Timeout limit for waiting for resources to finish applying.
`--verbose`&#124;`-V` | false | Prints underlying commands being called to stdout.

## Run Binary Locally

This assumes [gcloud](https://cloud.google.com/sdk/gcloud/) and
[kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) are installed
in your local machine and you have an existing GKE cluster to deploy to.

```bash
# Install gke-deploy.
go get github.com/GoogleCloudPlatform/cloud-builders/gke-deploy

# Run binary.
cd REPO  # Repo containing k8s YAML(s)
gke-deploy apply -f CONFIGS -n NAMESPACE -c CLUSTER_NAME -l CLUSTER_LOCATION
```

### Example

```bash
# Apply only.
gke-deploy apply -f ./configs -n my-namespace -c my-cluster -l us-east1-b

# Execute prepare and apply, with an intermediary step in between (e.g., manually check modified YAMLs).
gke-deploy prepare -f ./configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o ./modified -n my-namespace
cat modified/*
gke-deploy apply -f ./modified -n my-namespace -c my-cluster -l us-east1-b  # Pass modified directory to -f
```
