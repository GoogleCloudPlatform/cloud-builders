# Tool builder: `gcr.io/cloud-builders/helm`

This Container Builder build step runs
[`helm`](https://helm.sh/).

## Using this builder with Google Kubernetes Engine

To use this builder, it will need same permissions of [kubectl buidler](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/kubectl).

Please check [kubectl usage](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/kubectl#using-this-builder-with-google-kubernetes-engine).

Then, `helm` will have the configuration needed to talk to your GKE cluster.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud container builds submit . --config=cloudbuild.yaml
