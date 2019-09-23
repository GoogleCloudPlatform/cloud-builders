# Deploying with Cloud Build

This guide lists several simple use cases for using `gke-deploy` with Cloud
Build.

## Setup

These examples work with the following Google Cloud Platform APIs, which can
target different projects. **These examples assume a SINGLE PROJECT for all
APIs:**

* Cloud Build
* Container Registry
* Cloud Storage
* Kubernetes Engine

### Initialize project to run examples.

```bash
# Your project to run examples in.
PROJECT=my-project

# Create a cluster, if you have not done so.
CLUSTER=my-cluster
LOCATION=us-east1-b
gcloud container clusters create $CLUSTER --num-nodes=1 --zone=$LOCATION --project=$PROJECT

# Enable Cloud Build API, if you have not done so.
gcloud services enable cloudbuild.googleapis.com --project=$PROJECT

# Give the Cloud Build service account access to your project's clusters, if you have not done so.
PROJECT_NUM=$(gcloud projects list --filter="$PROJECT" --format="value(PROJECT_NUMBER)" --project=$PROJECT)
SERVICE_ACCOUNT=${PROJECT_NUM}@cloudbuild.gserviceaccount.com
gcloud projects add-iam-policy-binding $PROJECT --member=serviceAccount:$SERVICE_ACCOUNT --role=roles/container.developer --project=$PROJECT

# Create a bucket that will be used to store configs suggested and expanded by gke-deploy.
BUCKET=my-bucket
gsutil mb -p $PROJECT gs://$BUCKET
```

## Examples

### Build, publish, and deploy application with no Kubernetes configuration files

This build calls the `docker` build step to create a Docker image and publish it to
Container Registry. Then, the build calls the `gke-deploy` build step with the
`prepare` arg to create and expand Deployment, Horizontal Pod Autoscaler, and Service
configuration files to deploy. Next, the build calls the `gsutil` build step to
copy suggested base and expanded Kubernetes configuration files to your bucket.
Finally, the build calls the `gke-deploy` build step with the `apply` arg to
deploy the expanded Kubernetes configuration files created by the `prepare` step.

```bash
# Go to directory containing test app.
git clone https://github.com/GoogleCloudPlatform/cloud-builders.git && cd cloud-builders/gke-deploy/doc/app

# Variables for resources that were initialized above the Setup section.
PROJECT=my-project
BUCKET=my-bucket
CLUSTER=my-cluster
LOCATION=us-east1-b

# Arbitrary values
VERSION=my-version
APP=my-app
NAMESPACE=my-namespace

# Run build, replacing substitution variables accordingly.
gcloud builds submit . --project=$PROJECT --config cloudbuild-no-configs.yaml --substitutions=_IMAGE_NAME=gcr.io/$PROJECT/$APP,_IMAGE_VERSION=$VERSION,_GKE_CLUSTER=$CLUSTER,_GKE_LOCATION=$LOCATION,_K8S_APP_NAME=$APP,_K8S_NAMESPACE=$NAMESPACE,_OUTPUT_BUCKET_PATH=$BUCKET
```

### Build, publish, and deploy application with Kubernetes configuration files

This build calls the `docker` build step to create a Docker image and publish it to
Container Registry. Then, the build calls the `gke-deploy` build step with the
`prepare` arg to expand the provided configuration files to deploy. Next, the
build calls the `gsutil` build step to copy suggested base and expanded
Kubernetes configuration files to your bucket. Finally, the build calls the
`gke-deploy` build step with the `apply` arg to deploy the expanded Kubernetes
configuration files created by the `prepare` step.

```bash
# Go to directory containing test app.
git clone https://github.com/GoogleCloudPlatform/cloud-builders.git && cd cloud-builders/gke-deploy/doc/app

# Variables for resources that were initialized above the Setup section.
PROJECT=my-project
BUCKET=my-bucket
CLUSTER=my-cluster
LOCATION=us-east1-b

# Arbitrary values
VERSION=my-version
APP=my-app
NAMESPACE=my-namespace

# Replace image in config/deployment.yaml
sed -i "s#@IMAGE_NAME@#$IMAGE_NAME#g" config/deployment.yaml

# Replace namespace in config/namespace.yaml
sed -i "s#@NAMESPACE@#$NAMESPACE#g" config/namespace.yaml

# Run build, replacing substitution variables accordingly.
gcloud builds submit . --project=$PROJECT --config cloudbuild-with-configs.yaml --substitutions=_IMAGE_NAME=gcr.io/$PROJECT/$APP,_IMAGE_VERSION=$VERSION,_GKE_CLUSTER=$CLUSTER,_GKE_LOCATION=$LOCATION,_K8S_YAML_PATH=config,_K8S_APP_NAME=$APP,_K8S_NAMESPACE=$NAMESPACE,_OUTPUT_BUCKET_PATH=$BUCKET
```

You can remove the `Save configs` build step in your Cloud Build config if you do not
want to store configs to a bucket.

Builds can use `$SHORT_SHA` instead of an explicit `$_IMAGE_VERSION` for the
image and app versions if executed via a trigger. Follow [these
instructions](automated-deployments.md) to set up a trigger.
