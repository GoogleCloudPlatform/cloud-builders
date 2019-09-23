# Automated Deployments with Cloud Build

You can set up automated deployments of a repository in
[CSR](https://cloud.google.com/source-repositories/) using the following
command, which will cause a build, publish, and deploy to be triggered for every
push to any branch. This command allows setting up a build trigger without
needing an explicit `cloudbuild.yaml` file in your repository.

```bash
PROJECT=my-project  # Set project here
curl -X POST -H "Authorization: Bearer $(gcloud auth print-access-token --project=$PROJECT)" -H "Content-Type: application/json" -H "Accept: application/json" -H "X-Goog-User-Project: $PROJECT" https://cloudbuild.googleapis.com/v1/projects/$PROJECT/triggers -d '{
  "description": "Push to any branch",
  "triggerTemplate": {
    "repoName": "test-app",
    "branchName": ".*",
  },
  "build": {
    "steps": [
      {
        "name": "gcr.io/cloud-builders/docker",
        "id": "Build",
        "args": [
          "build",
          "-t",
          "$_IMAGE_NAME:$SHORT_SHA",
          ".",
          "-f",
          "$_DOCKERFILE_PATH",
        ]
      },
      {
        "name": "gcr.io/cloud-builders/docker",
        "id": "Publish",
        "args": [
          "push",
          "$_IMAGE_NAME:$SHORT_SHA",
        ]
      },
      {
        "name": "gcr.io/cloud-builders/gke-deploy:stable",
        "id": "Prepare deploy",
        "args": [
          "run",
          "--filename=$_K8S_YAML_PATH",
          "--image=$_IMAGE_NAME:$SHORT_SHA",
          "--app=$_K8S_APP_NAME",
          "--version=$SHORT_SHA",
          "--namespace=$_K8S_NAMESPACE",
          "--output=output",
          "--annotation=gcb-build-id=$BUILD_ID",
        ]
      },
      {
        "name": "gcr.io/cloud-builders/gsutil",
        "id": "Save configs",
        "entrypoint": "sh",
        "args": [
          "-c",
          "
set -e
gsutil cp -r output/suggested gs://$_OUTPUT_BUCKET_PATH/config/$BUILD_ID/suggested
echo \"Copied suggested base configs to gs://$_OUTPUT_BUCKET_PATH/config/$BUILD_ID/suggested\"
gsutil cp -r output/expanded gs://$_OUTPUT_BUCKET_PATH/config/$BUILD_ID/expanded
echo \"Copied expanded configs to gs://$_OUTPUT_BUCKET_PATH/config/$BUILD_ID/expanded\"
          "
        ]
      },
      {
        "name": "gcr.io/cloud-builders/gke-deploy:stable",
        "id": "Apply deploy",
        "args": [
          "run",
          "--filename=output/expanded",
          "--namespace=$_K8S_NAMESPACE",
          "--cluster=$_GKE_CLUSTER",
          "--location=$_GKE_LOCATION",
        ]
      }
    ],
    "images": [
      "$_IMAGE_NAME:$SHORT_SHA"
    ],
    "substitutions": {
      "_DOCKERFILE_PATH": "Dockerfile",
      "_IMAGE_NAME": "",
      "_GKE_CLUSTER": "",
      "_GKE_LOCATION": "",
      "_K8S_YAML_PATH": "",
      "_K8S_APP_NAME": "",
      "_K8S_NAMESPACE": "",
      "_OUTPUT_BUCKET_PATH": ""
    },
    "options": {
      "substitution_option": "ALLOW_LOOSE"
    },
    "tags": [
      "gcp-cloud-build-deploy",
      "$_K8S_APP_NAME"
    ]
  },
  "substitutions": {
    "_GKE_CLUSTER": "CLUSTER",
    "_GKE_LOCATION": "ZONE|REGION",
    "_IMAGE_NAME": "gcr.io/PROJECT/NAME",
    "_K8S_APP_NAME": "NAME",
    "_K8S_NAMESPACE": "NAMESPACE",
    "_K8S_YAML_PATH": "CONFIG",
    "_OUTPUT_BUCKET_PATH": "BUCKET_PATH"
  }
}'
```
