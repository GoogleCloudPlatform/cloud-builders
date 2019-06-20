# Bazel

This is a tool builder to simply invoke [`bazel`](https://bazel.io) commands.

You should consider instead using an [official `bazel`
image](https://launcher.gcr.io/google/bazel):

```yaml
steps:
- name: launcher.gcr.io/google/bazel:0.26.0
  args: ['version']
```

This allows you to use any supported version of Bazel.

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

If the rule is a [`docker_build`](https://bazel.build/versions/master/docs/be/docker.html#docker_build)
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

See https://bazel.build/versions/master/docs/be/docker.html for more options.

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
