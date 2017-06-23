# Docker

This is a tool builder to simply invoke `docker` commands.

Arguments passed to this builder will be passed to `docker` directly, allowing
callers to run [any Docker
command](https://docs.docker.com/engine/reference/commandline/).

By default, the version of Docker that is used by this builder is `17.05`.

## GCR Credentials

The Docker build step is automatically set up with credentials for your
[Container Builder Service
Account](https://cloud.google.com/container-builder/docs/permissions). These
permissions are sufficient to interact directly with GCR.

## Examples

The following examples demonstrate build requests that use this builder:

### Build and push a container image

This `cloudbuild.yaml` simply invokes `docker build` and pushes the resulting
image.

```
steps:
- name: gcr.io/cloud-builders/docker
  args: ['build', '-t', 'gcr.io/$PROJECT_ID/myimage', '.']
images: ['gcr.io/$PROJECT_ID/myimage']
```

### Retag an image and push the new tag

This `cloudbuild.yaml` adds a new tag (`:newtag`) to an existing image
(`:oldtag`) and pushes that tag.

```
steps:
- name: gcr.io/cloud-builders/docker
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
- name: gcr.io/cloud-builders/docker
  args: ['run', 'gcr.io/$PROJECT_ID/myimage']
```

## Older Versions

Since Docker CLI changes may not be backward-compatible, we provide tagged
versions of this builder for these previously-supported versions:

*   `gcr.io/cloud-builders/docker:1.12.6`
*   `gcr.io/cloud-builders/docker:17.05` (`:latest`)
