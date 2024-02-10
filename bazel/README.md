# Bazel

The `gcr.io/cloud-builders/bazel` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of Bazel. We also do
not provide historical pinned versions of bazel.

The Bazel team provides a `bazel` image that supports multiple tagged versions
at http://gcr.io/cloud-marketplace-containers/google/bazel.

To migrate to the Bazel team's official Bazel image, make the following changes
to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/bazel'
+ name: 'gcr.io/cloud-marketplace-containers/google/bazel'
+ entrypoint: 'bazel'
```

## Example Usage

```
steps:
- name: 'gcr.io/cloud-marketplace-containers/google/bazel'
  entrypoint: 'bazel'
  args: ['build', '//java/com/company/service:server']
```

---

## Usage Details

This is a tool builder to simply invoke [`bazel`](https://bazel.io) commands.

Arguments passed to this builder will be passed to `bazel` directly.

The latest available version of `bazel` is used.

## Examples

The following examples demonstrate build request that use this builder:

### Build a target

This `cloudbuild.yaml` simply builds a target. It might be a binary, a library,
a test target, or any other buildable target.

```
steps:
- name: gcr.io/cloud-builders/bazel
  args: ['build', '//java/com/company/service:server']
```

# Build and push a container image

If the rule is a [`docker_build`](https://docs.bazel.build/versions/0.17.1/be/docker.html#docker_build)
target, then you can `bazel run` the target to build a Docker image and load
it into the Docker daemon.  You can then tag the resulting image so it can be
pushed to the Container Registry.

`path/to/some/BUILD` file:

```
docker_build(
  name = "docker_target",
  base = "@docker_debian//:wheezy",
  entrypoint = ["echo", "foo"],
)
```

This `docker_build` rule produces a Docker container image based on debian that
specifies "echo foo" as its `ENTRYPOINT`.

See https://docs.bazel.build/versions/0.17.1/be/docker.html for more options.

`cloudbuild.yaml`:

```
steps:
# Build the Docker image and load it into the Docker daemon.
# The loaded image name is the BUILD target's name, prefixed with bazel/.
- name: gcr.io/cloud-builders/bazel
  args: ['run', '//path/to/some:docker_target']
# Re-tag the image to something in your project's gcr.io repository.
- name: gcr.io/cloud-builders/docker
  args: ['tag', 'bazel/path/to/some:docker_target', 'gcr.io/$PROJECT_ID/server']
# Push the image.
images: ['gcr.io/$PROJECT_ID/server']
```
