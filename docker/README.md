# Docker

The `gcr.io/cloud-builders/docker` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of Docker. We also do
not provide historical pinned versions of Docker.

The Docker Team provides a `docker` image that supports multiple tagged versions
as well as additional Docker tooling. For details, please visit
https://hub.docker.com/_/docker.

To migrate to the Docker Team's official image, make the following changes to
your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/docker'
+ name: 'docker'
```

## Example

Usage:

```yaml
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', '...']
```

## GCR Credentials

The Docker build step is automatically set up with credentials for your
[Cloud Build Service Account](https://cloud.google.com/cloud-build/docs/permissions).
These permissions are sufficient to interact directly with GCR.

## Examples

The following examples demonstrate build requests that use this builder:

### Build and push a container image

This `cloudbuild.yaml` simply invokes `docker build` and pushes the resulting
image.

```
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', 'gcr.io/$PROJECT_ID/myimage', '.']
images: ['gcr.io/$PROJECT_ID/myimage']
```

### Retag an image and push the new tag

This `cloudbuild.yaml` adds a new tag (`:newtag`) to an existing image
(`:oldtag`) and pushes that tag.

```
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['tag',
    'gcr.io/$PROJECT_ID/myimage:oldtag',
    'gcr.io/$PROJECT_ID/myimage:newtag']
images: ['gcr.io/$PROJECT_ID/myimage:newtag']
```

### Run a Docker image

This `cloudbuild.yaml` runs a Docker image. The step will complete when the
container exits. If the image runs too long, the build may time out.

```
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['run', 'gcr.io/$PROJECT_ID/myimage']
```
