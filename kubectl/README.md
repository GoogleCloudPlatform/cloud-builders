# Tool builder: `gcr.io/cloud-builders/kubectl`

This Container Builder build step runs
[`kubectl`](https://kubernetes.io/docs/user-guide/kubectl-overview/).

## Using this builder with Google Container Engine

To use this builder, your
[builder service account](https://cloud.google.com/container-builder/docs/how-to/service-account-permissions)
will need IAM permissions sufficient for the operations you want to perform. For
typical read-only usage, the "Container Engine Viewer" role is sufficient. To
deploy container images on a GKE cluster, the "Container Engine Developer" role
is sufficient. Check the
[GKE IAM page](https://cloud.google.com/container-engine/docs/iam-integration)
for details.

For most use, kubectl will need to be configured to point to a specific GKE
cluster. You can configure the cluster by setting environment variables.

    CLOUDSDK_COMPUTE_ZONE=<your cluster's zone>
    CLOUDSDK_CONTAINER_CLUSTER=<your cluster's name>

Setting the environment variables above will cause this step's entrypoint to
first run a command to fetch cluster credentials as follows.

    gcloud container clusters get-credentials --zone "$CLOUDSDK_COMPUTE_ZONE" "$CLOUDSDK_CONTAINER_CLUSTER"`

Then, `kubectl` will have the configuration needed to talk to your GKE cluster.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud container builds submit . --config=cloudbuild.yaml
