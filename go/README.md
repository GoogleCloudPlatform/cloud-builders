# Tool builder: `gcr.io/cloud-builders/go`

This builder runs the `go` tool (`go build`, `go test`, etc.)
after placing source in `/workspace` into the `GOPATH` before
running the tool.

# Using `golang` and [Go Modules](https://github.com/golang/go/wiki/Modules)

This Builder (`gcr.io/cloud-builders/go`) is not necessary if you're building using
[Go modules](https://github.com/golang/go/wiki/Modules), available
in Go 1.11+. You can **build** with the `golang` image (not `gcr.io/cloud-builders/go`) from [Dockerhub](https://hub.docker.com/_/golang) and Google's Container Registry [mirror](https://cloud.google.com/container-registry/docs/using-dockerhub-mirroring):

```
steps:
# If you already have a go.mod file, you can skip this step.
- name: mirror.gcr.io/library/golang
  args: ['go', 'mod', 'init', 'github.com/your/import/path']

# Build the module.
- name: mirror.gcr.io/library/golang
  env: ['GO111MODULE=on']
  args: ['go', 'build', './...']
```

See [`examples/module`](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/go/examples/module)
for a working example.

## Note #1 `/workspace` and `/go`
The `Golang` image defaults to a working directory of `/go` whereas Cloud Build mounts your sources under `/workspace`. This reflects the recommend best practice when using Modules of placing your sources *outside* of `GOPATH`. When you `go build ./...` you will run this in the working directory of `/workspace` *but* the packages will be pulled into `/go/pkg`.

Because `/go` is outside of `/workspace`, the `/go` directory is not persisted across Cloud Build steps. See next note.

## Note #2 Sharing packages across steps

One advantage with Go Modules is that packages are now semantically versioned and immutable; one a package has been pulled once, it should not need to be pulled again. Because the `golang` image uses `/go` as its working directory and this is outside of Cloud Build's `/workspace` directory, `/go` is recreated in each `Golang` step. To avoid this and share packages across steps, you may use Cloud Build `volumes`. An example to prove the point:

```YAML
- name: mirror.gcr.io/golang
  env:
  - GO111MODULE=on
  args: ['go','get','-u','github.com/golang/glog']
  volumes:
  - name: go-modules
    path: /go

- name: golang
  env:
  - GO111MODULE=on
  args: ['go','list','-f','{{ .Dir }}','-m','github.com/golang/glog']
  volumes:
  - name: go-modules
    path: /go
```
In the above, if the `volumes` section were omitted, the second step would fail. This is because `/go` would be created anew for the step and `github.com/golang/glob` would not be present in it.

**NB** Cloud Build supports using build-wide settings for `env` and `volumes` using `options` (see [link](https://cloud.google.com/cloud-build/docs/build-config#options)). I've duplicated here to aid clarity.

See [`examples/multi_step`](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/go/examples/multi_step/README.md)
for a working example.

## Note #3 Golang Module Mirror

The Go team provides a Golang Module Mirror ([https://proxy.golang.org/](https://proxy.golang.org/)). You may utilize the Mirror by including `GOPROXY=https://proxy.golang.org` in your build steps, e.g.:
```YAML
- name: mirror.gcr.io/golang
  env:
  - GO111MODULE=on
  - GOPROXY=https://proxy.golang.org
  args: ['go','get','-u','github.com/golang/glog']
  volumes:
  - name: go-modules
    path: /go

```

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

    $ gcloud builds submit . --config=cloudbuild.yaml
