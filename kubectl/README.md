# Tool builder: `gcr.io/cloud-builders/kubectl`

This Cloud Build build step runs
[`kubectl`](https://kubernetes.io/docs/user-guide/kubectl-overview/).

The `gcr.io/cloud-builders/kubectl` image is maintained by the Cloud Build team,
but it may not support the most recent versions of `kubectl`. We also do not
provide historical pinned versions of `kubectl` nor support across multiple
platforms.

The Cloud SDK team maintains a `cloud-sdk` image that supports multiple tagged
versions of Cloud SDK across multiple platforms. Please visit
https://github.com/GoogleCloudPlatform/cloud-sdk-docker for details.

Suggested alternative images include:

    gcr.io/google.com/cloudsdktool/cloud-sdk
    google/cloud-sdk

Please note that the `kubectl` entrypoint must be specified to use these images.

When executed in the Cloud Build environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/cloud-build/docs/permissions) for the build
project.

To migrate to the Cloud SDK team's official image, make the following changes
to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/kubectl'
+ name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
+ entrypoint: 'kubectl'
```

-----

## Usage

For most uses, `kubectl` will need to be configured to point to a specific GKE
cluster. You can configure the cluster by setting environment variables.

    # Set region for regional GKE clusters or Zone for Zonal clusters
    CLOUDSDK_COMPUTE_REGION=<your cluster's region>
    # or
    CLOUDSDK_COMPUTE_ZONE=<your cluster's zone>

    # Name of GKE cluster
    CLOUDSDK_CONTAINER_CLUSTER=<your cluster's name>

**When using Google Cloud Build, you must set these environment variables on
every step that uses the `kubectl` builder; this context is not persisted across
steps.**

If your GKE cluster is in a different project than the build itself, also set:

```CLOUDSDK_CORE_PROJECT=<the GKE cluster project>```

## Using this builder with Google Kubernetes Engine

To use this builder on Google Cloud Build, your [builder service
account](https://cloud.google.com/cloud-build/docs/how-to/service-account-permissions)
will need IAM permissions sufficient for the operations you want to perform.

For typical read-only usage, enable the "Kubernetes Engine Viewer" role. Check the
[GKE IAM page](https://cloud.google.com/kubernetes-engine/docs/how-to/iam#roles) for
details.

To deploy container images on a GKE cluster, enable the "Kubernetes Engine Developer"
role:

- Open the [Cloud Build settings page](https://console.cloud.google.com/cloud-build/settings).
- You'll see the **Service account permissions** page.
- Set the status of the "Kubernetes Engine Developer" role to **Enable**.

Make sure you also grant the Cloud Build service account permissions in the GKE cluster project.

---

## Additional Usage Information

When using `gcr.io/cloud-builders/kubectl`, setting the environment variables
described above will cause this step's entrypoint to first run a command to
fetch cluster credentials as follows.

    gcloud container clusters get-credentials --zone "$CLOUDSDK_COMPUTE_ZONE" "$CLOUDSDK_CONTAINER_CLUSTER"

Then, `kubectl` will have the configuration needed to talk to your GKE cluster.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
