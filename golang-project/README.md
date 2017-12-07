# !!! `golang-project` is deprecated !!!

**Use `gcr.io/cloud-builders/go` instead**

## Migrating to `gcr.io/cloud-builders/go`

When you run this builder, you should see some log lines like:

```
Dockerfile contents:
--------------------
FROM <BASE IMAGE>
ENV_PATH=/golang_project_bin:$PATH
COPY gopath/bin/hello /golang_project_bin/
ENTRYPOINT ["/golang_project_bin/<PACKAGE>"]
--------------------
```

Copy these contents into a new file named `Dockerfile` and fill in your
preferred base image and binary name.

Update your `cloudbuild.yaml` config. Instead of something like this:

```
steps:
- name: 'gcr.io/cloud-builders/golang-project'
  args: [<PACKAGE>, <IMAGE>]
images: [<IMAGE>]
```

Use the `go` builder:

```
steps:
- name: 'gcr.io/cloud-builders/go'
  args: ['install']
- name: 'gcr.io/cloud-builders/go'
  args: ['test', './...']
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', <IMAGE>, '.']
images: [<IMAGE>]
```

This build configuration will build and test your Go package, then build a
Docker image using the new `Dockerfile`.

See [hello_world_migrated](examples/hello_world_migrated) for an example of
the [hello_world example](examples/hello_world), migrated to use the `go`
builder.

# Native builder: `gcr.io/cloud-builders/golang-project`

This Container Builder build step builds canonical Go projects.

### When to use this builder

The `gcr.io/cloud-builders/golang-project` build step should be used when you
have a "typical" Go source tree and want a single step to turn it into a
container image. This step relies on conventions described in the "Workspace
setup" section and shown in the [examples](examples).

The process is:

1.  Test and build all the targets indicated by the step's `args` using `go
    test` and `go install`.
1.  Create a `Dockerfile` on the fly, using the `--base-image` and
    `--entrypoint` provided.
1.  Build the container image using the new `Dockerfile` and the `--tag`
    provided.

To customize the image created, either use a `--base-image` with those
customizations, or use the [`gcr.io/cloud-builders/docker`](../docker/README.md)
build step with a `Dockerfile` that is built `FROM` the image produced by
the `gcr.io/cloud-builders/golang-project` build step.

### `alpine` vs `wheezy`

There are two versions of this builder, one for `alpine` and one for `wheezy`.
The difference is significant since, by default, Go dynamically links libc.
Binaries built in an `alpine` environment don't always work in a `wheezy`-like
(including `ubuntu`, etc) environment, or vice versa.

The specific versions are available as

  - gcr.io/cloud-builders/golang-project:alpine
  - gcr.io/cloud-builders/golang-project:wheezy

And `gcr.io/cloud-builders/golang-project:latest` is an alias for
`gcr.io/cloud-builders/golang-project:alpine`.

### Related: `gcr.io/cloud-builders/go`

The related build step, [`gcr.io/cloud-builders/go`](../go/README.md), is used
to run the `go` directly on source. The [`gcr.io/cloud-builders/go`]
(../go/README.md) build step and this `gcr.io/cloud-builders/golang-project`
recognize Go workspaces in the same way, described in the
[`gcr.io/cloud-builders/go`](../go/README.md) build step's [README]
(../go/README.md)

## Workspace setup

The native Go builder uses the same workspace setup mechanisms as
`gcr.io/cloud-builders/go`. See its [documentation](../go/README.md) for
details.

## Usage

The positional arguments are the targets that will be built using `go install`,
and copied into the resulting container image.

### `--base-image=BASE_IMAGE`

This flag chooses the base image (default is `alpine`) used when building the
resulting container image.

### `--tag=TAG`

This flag, which is required, chooses the tag used when building the resulting
container image.

### `--skip-tests`

If this flag is set, `go test` is not run.

### `--entrypoint=ENTRYPOINT`

This flag indicates what entrypoint should be used in the resulting container
image. If not set, the first executable target is used.

## Examples

-   [Hello, World!](examples/hello_world) is a basic example that creates a
    container using the indicated target as its entrypoint.
-   [Multiple targets](examples/multi_bin) creates two runnable targets, and
    shows how to indicate which one is the entrypoint.
-   [Static webserver](examples/static_webserver) creates a container image that
    combines a Go project with static resources, and updates the image's
    entrypoint to indicate those resources.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud container builds submit . --config=cloudbuild.yaml
