# Automated Deployments with Cloud Build

You can set up automated deployments of a GitHub repository using the following
command, which will cause a build, publish, and deploy to be triggered for every
push to any branch. This command allows setting up a build trigger without
needing an explicit `cloudbuild.yaml` file in your repository.

```bash
# Go to directory containing this doc.
git clone https://github.com/GoogleCloudPlatform/cloud-builders.git && cd cloud-builders/gke-deploy/doc

PROJECT=my-project
OWNER=user
REPO=repo
BRANCH_REGEX='.*'
IMAGE_NAME=image
CLUSTER=cluster
LOCATION='zone|region'
CONFIGS=configs
APP_NAME=name
NAMESPACE=namespace
OUTPUT_BUCKET_PATH=gs://bucket/path

# Copy trigger.yaml and replace values
cp trigger.yaml my-trigger.yaml
sed -i "s#@OWNER@#$OWNER#g" my-trigger.yaml
sed -i "s#@REPO@#$REPO#g" my-trigger.yaml
sed -i "s#@BRANCH_REGEX@#$BRANCH_REGEX#g" my-trigger.yaml
sed -i "s#@IMAGE_NAME@#$IMAGE_NAME#g" my-trigger.yaml
sed -i "s#@CLUSTER@#$CLUSTER#g" my-trigger.yaml
sed -i "s#@LOCATION@#$LOCATION#g" my-trigger.yaml
sed -i "s#@CONFIGS@#$CONFIGS#g" my-trigger.yaml
sed -i "s#@APP_NAME@#$APP_NAME#g" my-trigger.yaml
sed -i "s#@NAMESPACE@#$NAMESPACE#g" my-trigger.yaml
sed -i "s#@OUTPUT_BUCKET_PATH@#$OUTPUT_BUCKET_PATH#g" my-trigger.yaml

gcloud alpha builds trigger create github --trigger-config=my-trigger.yaml --project=$PROJECT
```
