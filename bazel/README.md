# Bazel

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

If the rule is a [`docker_build`](https://www.bazel.io/docs/be/docker.html)
target, then you can `bazel run` that image to produce a container image, then
push that image.

`path/to/some/BUILD` file:

```
docker_build(
  name = "docker_target",
  base = "@docker_debian//:wheezy"
  entrypoint = ["echo", "foo"],
)
```

This `docker_build` rule produces a Docker container image based on debian that
specifies "echo foo" as its ENTRYPOINT.

See https://www.bazel.io/docs/be/docker.html for more options.

`cloudbuild.yaml`:

```
steps:
# Build the Docker image. The result is an image based the target's name.
- name: gcr.io/cloud-builders/bazel
  args: ['run', '//path/to/some:docker_target']
# Re-tag the image to something in your project's gcr.io repository.
  name: gcr.io/cloud-builders/docker
  args: ['tag', 'bazel/path_to_some:docker_target', 'gcr.io/$PROJECT_ID/server']
# Push the image.
images: ['gcr.io/$PROJECT_ID/server']
```
