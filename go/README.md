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
is replaced by the use of [Go modules](https://golang.org/ref/mod), available
since Go 1.11. As of Go 1.16, go modules are the default behavior.

See [`examples/module`](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/go/examples/module)
for a working example.

## Note #1 `/workspace` and `/go`
The `Golang` image defaults to a working directory of `/go` whereas Cloud Build
mounts your sources under `/workspace`. This reflects the recommended best
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
```
Note that this examples defines `volumes` as a global build option (rather than
on each build step). This allows the `/go` directory to be shared across steps
-- it's initialized in the first and used in the second.

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
  - GOPROXY=https://proxy.golang.org
  args: ['go','get','-u','github.com/golang/glog']
  volumes:
  - name: go-modules
    path: /go
```
----

## Using `gcr.io/cloud-builders/go`

This builder runs the `go` tool (`go build`, `go test`, etc.)
after placing source in `/workspace` into the `GOPATH` before
running the tool.

### `alpine` vs `debian`

There are two supported versions of this builder, one for `alpine` and one for
`debian`. The difference is significant since, by default, Go dynamically links
`libc`. Binaries built in an `alpine` environment don't always work in a
`debian`-like (including `ubuntu`, etc) environment, or vice versa.

The specific versions are available as

  - gcr.io/cloud-builders/go:alpine
  - gcr.io/cloud-builders/go:debian

And `gcr.io/cloud-builders/go:latest` is an alias for
`gcr.io/cloud-builders/go:alpine`.

This tool builder sets `CGO_ENABLED=0`, so that all binaries are linked
statically.

## Output files

If you use the `install` subcommand, the binaries will end up in `$GOPATH/bin`
which is unlikely to be in your `/workspace`. You may want to use `go build -o
<target>` instead.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
