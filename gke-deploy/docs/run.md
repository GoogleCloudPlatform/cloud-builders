# gke-deploy run [flags]

This command will execute both the [Prepare](../README.md#prepare-step) and
[Apply](../README.md#apply-step) phases.

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
gke-deploy run -f CONFIGS -i IMAGE -a APP_NAME -v APP_VERSION -o OUTPUT -n NAMESPACE -c CLUSTER_NAME -l CLUSTER_LOCATION
```

### Example

```bash
# Modify configs and deploy to GKE cluster.
gke-deploy run -f ./configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o ./modified -n my-namespace -c my-cluster -l us-east1-b

# Deploy to GKE cluster that kubectl is currently targeting.
gke-deploy run -f ./configs
```

## Run Docker Container Locally

This assumes [gcloud](https://cloud.google.com/sdk/gcloud/) and
[docker](https://docs.docker.com/) are installed in your local machine and you
have an existing GKE cluster to deploy to.

```bash
docker run -rm -v ~/.config/gcloud:/root/.config/gcloud -v ~/.gcp/root/.gcp -e GOOGLE_APPLICATION_CREDENTIALS=/root/.gcp/service-account.json -v /absolute/path/to/k8s/configs:/config gcr.io/cloud-builders/gke-deploy run -f /config -i image -a name -v version -n namespace -c cluster_name -l cluster_location
```

## Run on Google Cloud Build

This assumes [gcloud](https://cloud.google.com/sdk/gcloud/) is installed in your
local machine and you have an existing GKE cluster to deploy to. You can also
pass a bucket to the `_OUTPUT_BUCKET` substitution variable to save the modified
YAMLs.

```bash
# Enable Cloud Build API, if you have not done so.
gcloud services enable cloudbuild.googleapis.com

# Give the Cloud Build service account access to your clusters, if you have not done so.
PROJECT=$(gcloud config get-value project) && \
PROJECT_NUM=$(gcloud projects list --filter="$PROJECT" --format="value(PROJECT_NUMBER)") && \
gcloud projects add-iam-policy-binding $PROJECT --member=serviceAccount:${PROJECT_NUM}@cloudbuild.gserviceaccount.com --role=roles/container.developer

# Create cloudbuild.yaml
cd repo  # Repo containing k8s YAML(s)
cat >cloudbuild.yaml <<"EOF"
steps:
- name: 'gcr.io/cloud-builders/docker'
  id: 'Build'
  args:
  - 'build'
  - '-t'
  - '$_IMAGE_NAME:$_IMAGE_VERSION'
  - '.'
  - '-f'
  - '$_DOCKERFILE_PATH'
- name: 'gcr.io/cloud-builders/docker'
  id: 'Push'
  args:
  - 'push'
  - '$_IMAGE_NAME:$_IMAGE_VERSION'
- name: 'gcr.io/cloud-builders/gke-deploy'
  id: 'Deploy'
  args:
  - 'run'
  - '--filename=$_K8S_YAML_PATH'
  - '--image=$_IMAGE_NAME:$_IMAGE_VERSION'
  - '--app=$_K8S_APP_NAME'
  - '--version=$_IMAGE_VERSION'
  - '--namespace=$_K8S_NAMESPACE'
  - '--label=$_K8S_LABELS,gcb-build-id=$BUILD_ID'
  - '--output=output'
  - '--cluster=$_GKE_CLUSTER'
  - '--location=$_GKE_LOCATION'
images:
- '$_IMAGE_NAME:$_IMAGE_VERSION'
artifacts:
  objects:
    location: 'gs://$_OUTPUT_BUCKET/$BUILD_ID'
    paths: ['output/*']
substitutions:
  _DOCKERFILE_PATH: Dockerfile
  _IMAGE_NAME:
  _IMAGE_VERSION:
  _GKE_CLUSTER:
  _GKE_LOCATION:
  _K8S_YAML_PATH:
  _K8S_APP_NAME:
  _K8S_NAMESPACE:
  _K8S_LABELS:
  _OUTPUT_BUCKET:
options:
  substitution_option: 'ALLOW_LOOSE'
tags: ['gcp-cloud-build-deploy', '$_K8S_APP_NAME']
EOF

# Run build
gcloud builds submit . --config cloudbuild.yaml --substitutions=_IMAGE_NAME=gcr.io/project/name,_IMAGE_VERSION=version,_GKE_CLUSTER=cluster,_GKE_LOCATION=zone|region,_K8S_YAML_PATH=config,_K8S_APP_NAME=name,_K8S_NAMESPACE=namespace,_OUTPUT_BUCKET=bucket

# You can remove the artifacts field if you do not want to store the output into a bucket.

# Builds can use `$SHORT_SHA` instead of an explicit `$_IMAGE_VERSION` for the image and app versions if executed via a trigger.
```
