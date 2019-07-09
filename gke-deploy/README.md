# GKE Deploy

** Warning: This cloud builder is experimental and will very likely change in
breaking ways at this time. **

This tool deploys an application to a GKE cluster, following Google's
recommended best practices.

## gke-deploy vs kubectl

Using `gke-deploy` to deploy an application to GKE differs from kubectl in that
a `gke-deploy` deployment follows Google's recommended best practices by doing
the following:

### Prepare Step

* The `gke-deploy` builder modifies a set of Kubernetes resource YAML configs to use a
container image's digest instead of a tag.
* This builder adds [recommended
labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#applications-and-instances-of-applications).

### Apply Step

The `gke-deploy`builder does the following:

1. Gets authorized to access a GKE cluster

2. Applies the set of Kubernetes resource YAML configs that were
modified in the prepare step

3. Waits for applied Kubernetes resources to be ready

## Usage

`[gke-deploy run [flags]](docs/run.md)`

This command will execute both the [Prepare](#prepare-step) and
[Apply](#apply-step) phases mentioned above.

`[gke-deploy prepare [flags]](docs/prepare.md)`

This command will execute only the [Prepare](#prepare-step) phase mentioned
above.

`[gke-deploy apply [flags]](docs/apply.md)`

This command will execute only the [Apply](#apply-step) phase mentioned above.

### [Automated Deployments with GCB](docs/automated-deployments.md)

Follow [these instructions](docs/automated-deployments.md) to set up continuous deployment.
