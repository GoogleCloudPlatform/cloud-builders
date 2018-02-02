# Tool builder: gcr.io/cloud-builders/dotnet

This directory contains the description for the `dotnet` Container Build step.
This step will run the .NET Core `dotnet` command line. This builder has two
SDKs installed:

+ 1.0.3 with SDK Preview 2 build 3156, to support .NET Core 1.0.5 with
  project.json solutions.
+ 1.1.7, which supports .csproj projects for .NET Core 1.0.9 and 1.1.6.
+ 2.1.4, which supports projects for .NET Core 2.0.5

Any arguments passed to this builder step will be passed directly to the
`dotnet` tool. Note that when using the `dotnet restore` command you should
always specify the `--packages` option to a path within the `/workspace`, this
way the restored packages will be available to the next invocations of `dotnet`.
See the [example](examples/TestApp/cloudbuild.yaml) for how use `dotnet restore`
correctly.

We strongly advice that all .NET Core projects contain a `global.json` file that
specifies what version of the SDK is to be used. Otherwise `dotnet` will just
use the latest version of the SDK available. This is specially true for old
projects still using the `project.json` build system, you need make sure that
you specify the right SDK so an SDK that supports `project.json` is used to
build your project.

## Building this Container Build step

To build this builder, run the following command on this directory:

```bash
gcloud container builds submit --config=./cloudbuild.yaml .
```

This will build the Container Build step in your current GCP project.
