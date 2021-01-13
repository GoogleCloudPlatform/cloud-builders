# Tool builder: `gcr.io/cloud-builders/go`

The `gcr.io/cloud-builders/go` image is maintained by the Cloud Build team, but
it may not support the most recent features or versions of Go. We also do not
provide historical pinned versions of Go.

The Go team on Dockerhub provides the
[`golang`](https://hub.docker.com/_/golang) image which includes additional Go
tooling and supports multiple tagged versions across several platforms.

To migrate to the Go team's official `golang` image, make the following changes
to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/go'
+ name: 'golang'
+ entrypoint: 'go'
```

# Go Modules

The functionality previously provided by `gcr.io/cloud-builders/go` in
[`prepare_workspace`](https://github.com/GoogleCloudPlatform/cloud-builders/blob/master/go/prepare_workspace.inc)
is replaced by the use of [Go
modules](https://github.com/golang/go/wiki/Modules), available since Go 1.11.

## Example

```
steps:
# If you already have a go.mod file, you can skip this step.
- name: 'golang'
  args: ['go', 'mod', 'init', 'github.com/your/import/path']

# Build the module.
- name: 'golang'
  env: ['GO111MODULE=on']
  args: ['go', 'build', './...']
```

See [`examples/module`](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/go/examples/module)
for a working example.

## Note #1 `/workspace` and `/go`
The `Golang` image defaults to a working directory of `/go` whereas Cloud Build
mounts your sources under `/workspace`. This reflects the recommend best
practice when using Modules of placing your sources *outside* of `GOPATH`. When
you `go build ./...` you will run this in the working directory of `/workspace`
*but* the packages will be pulled into `/go/pkg`.

Because `/go` is outside of `/workspace`, the `/go` directory is not persisted
across Cloud Build steps. See next note.

## Note #2 Sharing packages across steps

One advantage with Go Modules is that packages are now semantically versioned
and immutable; once a package has been pulled, it should not need to be
pulled again. Because the `golang` image uses `/go` as its working directory and
this is outside of Cloud Build's `/workspace` directory, `/go` is recreated in
each `Golang` step. To avoid this and share packages across steps, you may use
Cloud Build `volumes`.

Here is an example to prove the point; this example uses `go mod init` as the
first step rather than creating a `go.mod`:

```YAML
steps:
- name: 'golang'
  args: ['go','mod','init','github.com/golang/glog']
- name: 'golang'
  args: ['go','list','-f','{{ .Dir }}','-m','github.com/golang/glog']
options:
  volumes:
  - name: go-modules
    path: /go
  env:
  - GO111MODULE=on
```
Note that this examples defines `volumes` and `env` as global build options
(rather than on each build step). This allows the `/go` directory to be shared
across stepped -- initialized in the first and used in the second.

**NB** Cloud Build supports using build-wide settings for `env` and `volumes`
using `options` (see
[link](https://cloud.google.com/cloud-build/docs/build-config#options)).

See [`examples/multi_step`](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/go/examples/multi_step/README.md)
for a working example.

## Note #3 Golang Module Mirror

The Go team provides a Golang Module Mirror
([https://proxy.golang.org/](https://proxy.golang.org/)). The module mirror is
enabled by default for Go 1.13 and newer. You may utilize the Mirror by
including `GOPROXY=https://proxy.golang.org` in your build steps, e.g.:
```YAML
- name: golang
  env:
  - GO111MODULE=on
  - GOPROXY=https://proxy.golang.org
  args: ['go','get','-u','github.com/golang/glog']
  volumes:
  - name: go-modules
    path: /go
```
----
## Historical Usage

This builder runs the `go` tool (`go build`, `go test`, etc.)
after placing source in `/workspace` into the `GOPATH` before
running the tool.

## Using `gcr.io/cloud-builders/go`

### `alpine` vs `debian`

There are two supported versions of this builder, one for `alpine` and one for
`debian`. The difference is significant since, by default, Go dynamically links
libc. Binaries built in an `alpine` environment don't always work in a
`debian`-like (including `ubuntu`, etc) environment, or vice versa.

The specific versions are available as

  - gcr.io/cloud-builders/go:alpine
  - gcr.io/cloud-builders/go:debian

And `gcr.io/cloud-builders/go:latest` is an alias for
`gcr.io/cloud-builders/go:alpine`.

## Workspace setup

Before the `go` tool is used, the build step first sets up a workspace.

To determine the workspace structure, this tool checks the following, in order:

1.  Is `$GOPATH` set? Use that.
2.  Is `$PROJECT_ROOT` set? Make a temporary workspace in `GOPATH=./gopath`, and
    link the contents of the current directory into
    `./gopath/src/$PROJECT_ROOT/*`.
3.  Is there a `./src` directory? Set `GOPATH=$PWD`.
4.  Does a `.go` file in the current directory have a comment like `// import
    "$PROJECT_ROOT"`? Use the `$PROJECT_ROOT` found in the import comment.

Once the workspace is set up, the `args` to the build step are passed through to
the `go` tool.

This tool builder sets `CGO_ENABLED=0`, so that all binaries are linked statically.

## Output files

The binaries built by a `go install` step will be available to subsequent build
steps.

If you use the `install` subcommand, the binaries will end up in `$GOPATH/bin`.

*   If you set `$GOPATH`, the binaries will end up in `$GOPATH/bin`.
*   If you rely on `./src` to indicate the workspace, the binaries will end up
    in `./bin`.
*   If you use `$PROJECT_ROOT` or an `//import` comment, the binaries will end
    up in `./gopath/bin`.

## Examples

-   [Hello, World!](examples/hello_world) is a basic example that builds a
    binary that is injected into a container image. It uses the `$PROJECT_ROOT`
    method to define its workspace.
-   [Whole workspace](examples/whole_workspace) is the same as the "Hello,
    World!" example, except that it uses the `./src` method to define its
    workspace.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
