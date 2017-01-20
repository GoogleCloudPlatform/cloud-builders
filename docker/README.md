# Docker

This is a tool builder to simply invoke `docker` commands.

Arguments passed to this builder will be passed to `docker` directly, allowing
callers to run [any Docker
command](https://docs.docker.com/engine/reference/commandline/).

The version of Docker that is used by this builder is `1.9.1`.

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
- name: gcr.io.cloud-builders/docker
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
