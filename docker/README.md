# Docker

## Deprecation Notice

This tool builder to invoke `docker` commands is deprecated in favor of the
[official `docker` images](https://hub.docker.com/_/docker/) on Dockerhub.

```yaml
steps:
- name: 'docker'
  args: ['build', '-t', '...']
```

Note tha the official images support many tagged version of Docker across many
base image configurations.

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
- name: 'docker'
  args: ['build', '-t', 'gcr.io/$PROJECT_ID/myimage', '.']
images: ['gcr.io/$PROJECT_ID/myimage']
```

### Retag an image and push the new tag

This `cloudbuild.yaml` adds a new tag (`:newtag`) to an existing image
(`:oldtag`) and pushes that tag.

```
steps:
- name: 'docker'
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
- name: 'docker'
  args: ['run', 'gcr.io/$PROJECT_ID/myimage']
```
