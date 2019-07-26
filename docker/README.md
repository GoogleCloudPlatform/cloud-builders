# Docker

This is a tool builder to simply invoke `docker` commands.

Arguments passed to this builder will be passed to `docker` directly, allowing
callers to run
[any Docker command](https://docs.docker.com/engine/reference/commandline/).

You should consider instead using an [official `docker`
image](https://hub.docker.com/_/docker/):

```yaml
steps:
- name: docker
  args: ['build', '-t', '...']
```

This allows you to use any supported version of Docker (e.g., `docker:19.03.1`)

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

Since Docker CLI changes may not be backward-compatible, tagged versions
of this builder remain available for several previously-supported versions:

*   `gcr.io/cloud-builders/docker:18.06.1`
*   `gcr.io/cloud-builders/docker:18.09.0`
*   `gcr.io/cloud-builders/docker:18.09.6` (`:latest`)
