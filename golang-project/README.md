# Native builder: `gcr.io/cloud-builders/golang-project`

This Container Builder build step builds canonical Go projects.

The process is:

1.  Test and build all the targets indicated by the step's `args` using
    `go test` and `go install`.
1.  Create a `Dockerfile` on the fly, using the `--base-image` and
    `--entrypoint` provided.
1.  Build the container image using the new `Dockerfile` and the `--tag`
	provided.

## Workspace setup

The native Go builder uses the same workspace setup mechanisms as
`gcr.io/cloud-builders/go`. See its [documentation](../go/README.md) for
details.

## Usage

The positional arguments are the targets that will be built using
`go install`, and copied into the resulting container image.

### `--base-image=BASE_IMAGE`

This flag chooses the base image (default is `alpine`) used when building the
resulting container image.

### `--tag=TAG`

This flag, which is required, chooses the tag used when building
the resulting container image.

### `--skip-tests`

If this flag is set, `go test` is not run.

### `--entrypoint=ENTRYPOINT`

This flag indicates what entrypoint should be used in the resulting container
image. If not set, the first executable target is used.

## Examples

- [Hello, World!](examples/hello_world) is a basic example that creates a
  container using the indicated target as its entrypoint.
- [Multiple targets](examples/multi_bin) creates two runnable targets, and
  shows how to indicate which one is the entrypoint.
- [Static webserver](examples/static_webserver) creates a container image
  that combines a Go project with static resources, and updates the image's
  entrypoint to indicate those resources.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud alpha container builds create . --config=cloudbuild.yaml
