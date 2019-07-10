# gke-deploy prepare [flags]

This command will execute only the [Prepare](../README.md#prepare-step) phase.

## Parameters

Flag | Default Value | Description
--- | --- | --- |
`--filename`&#124;`-f` | | Config file or directory of config files to use to create the Kubernetes resources (file or files in directory must end with \".yml\" or \".yaml\"). If this is omitted, base YAMLs will be generated for the user.
`--image` &#124;`-i` | | Image(s) to be deployed. Images can be set comma-delimited or as separate flags.
`--app` &#124;`-a` | | Name of the Kubernetes deployment.
`--version`&#124;`-v` | | Version of the Kubernetes deployment.
`--output` &#124;`-o` | `./output` | Target directory to store outputs. These files can be copied to a bucket for auditting purposes.
`--namespace`&#124;`-n` | `default` | Kubernetes namespace in which to deploy the app.
`--label`&#124;`-L` | | Label(s) to add to Kubernetes resources (k1=v1). Labels can be set comma-delimited or as separate flags. If two or more labels with the same key are listed, the last one is used.
`--verbose`&#124;`-V` | false | Prints underlying commands being called to stdout.

## Run Binary Locally

This assumes [gcloud](https://cloud.google.com/sdk/gcloud/) is installed in your
local machine.

```bash
# Install gke-deploy.
go get github.com/GoogleCloudPlatform/cloud-builders/gke-deploy

# Run binary.
cd REPO  # Repo containing k8s YAML(s)
gke-deploy prepare -f CONFIGS -i IMAGE -a APP_NAME -v APP_VERSION -o OUTPUT -n NAMESPACE
```

### Example

```bash
# Prepare only.
gke-deploy prepare -f ./configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o ./modified -n my-namespace

# Execute prepare and apply, with an intermediary step in between (e.g., manually check modified YAMLs).
gke-deploy prepare -f ./configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o ./modified -n my-namespace
cat modified/*
gke-deploy apply -f ./modified -c my-cluster -n my-namespace -c my-cluster -l us-east1-b  # Pass modified directory to -f
```
