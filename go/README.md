# Tool builder: `gcr.io/cloud-builders/go`

This Container Builder build step runs the `go` tool.

## Workspace setup

Before the `go` tool is used, the build step first sets up a workspace.

To determine the workspace structure, this tool checks the following, in order:

1. Is `$GOPATH` set? Use that.
2. Is there a `./src` directory? Set `GOPATH=$PWD`.
3. Is `$PROJECT_ROOT` set? Make a temporary workspace in `GOPATH=./gopath`, and link
   the contents of the current directory into `./gopath/src/$PROJECT_ROOT/*`.
4. Does a `.go` file in the current directory have a comment like
   `// import "$PROJECT_ROOT"`? Use the `$PROJECT_ROOT` found in the import
   comment instead of a provided `$PROJECT_ROOT` environment variable.

Once the workspace is set up, the `args` to the build step are passed through to the `go` tool.

## Output files

The binaries built by a `go install` step will be available to subsequent build
steps.

If you use the `install` subcommand, the binaries will end up in `$GOPATH/bin`.

* If you set `$GOPATH`, the binaries will end up in `$GOPATH/bin`.
* If you rely on `./src` to indicate the workspace, the binaries will end up
  in `./bin`.
* If you use `$PROJECT_ROOT` or an `//import` comment, the binaries will end up
  in `./gopath/bin`.

## Examples

- [Hello, World!](examples/hello_world) is a basic example that builds a binary
  that is injected into a container image. It uses the `$PROJECT_ROOT` method to
  define its workspace.
- [Whole workspace](examples/whole_workspace) is the same as the "Hello, World!"
  example, except that it uses the `./src` method to define its workspace.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud alpha container builds create . --config=cloudbuild.yaml
