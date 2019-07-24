# GKE Deploy

This tool deploys an application to a Kubernetes Engine cluster, following
Google's recommended best practices.

## Install Locally

Although `gke-deploy` is meant to be used as a build step with [Cloud
Build](https://cloud.google.com/cloud-build/), that doesn't mean that it can't
be used locally.

1.  First, install
    [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) as a
    dependency

2.  Next, install `gke-deploy`

    ```bash
    go install github.com/GoogleCloudPlatform/cloud-builders/gke-deploy
    gke-deploy -h
    ```

3.  If your `kubectl` is pointing to a cluster, you can test out a deployment by
    deploying an application with one simple command

    ```bash
    # Deploy an nginx Deployment with a load balancer exposing port 80.
    gke-deploy run -i nginx -x 80

    # After the command finishes, you should be able to visit the IP Address printed
    # in the > Deployed Resources table in the Service row.
    # e.g., Service                    nginx-service    Yes      35.231.198.229
    curl 35.231.198.229
    ```

## gke-deploy vs kubectl

Using `gke-deploy` to deploy an application to Kubernetes Engine differs from
`kubectl` in that `gke-deploy` is effectively a wrapper around a `kubectl apply`
deployment that follows Google's recommended best practices by doing the
following:

### Prepare Step

*   The `gke-deploy` builder modifies a set of Kubernetes resource YAML configs
    to use a container image's digest instead of a tag.

*   The builder adds several [recommended
    labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels).

### Apply Step

`gke-deploy` does the following:

1.  Gets authorized to access a GKE cluster

2.  Applies (using `kubectl`) the set of Kubernetes resource YAML configs that were
    modified in the prepare step

3.  Waits for applied Kubernetes resources to be ready

## Usage

[`gke-deploy run [flags]`](doc/gke-deploy_run.md)

This command will execute both the [Prepare](#prepare-step) and
[Apply](#apply-step) phases mentioned above.

[`gke-deploy prepare [flags]`](doc/gke-deploy_prepare.md)

This command will execute only the [Prepare](#prepare-step) phase mentioned
above.

[`gke-deploy apply [flags]`](doc/gke-deploy_apply.md)

This command will execute only the [Apply](#apply-step) phase mentioned above.

## [Deploying with Cloud Build](doc/deploying-with-cloud-build.md)

View [this page](doc/deploying-with-cloud-build.md) for examples on how to use
`gke-deploy` with Cloud Build.

## [Automated Deployments with Cloud Build](doc/automated-deployments.md)

Follow [these instructions](doc/automated-deployments.md) to set up continuous
deployment.
