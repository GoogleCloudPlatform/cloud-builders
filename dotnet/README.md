# Tool builder: gcr.io/cloud-builders/dotnet

This directory contains the description for the `dotnet` Container Build step.
This step will run the .NET Core `dotnet` command line. This builder has three
SDKs installed:

+ 1.0.3 with SDK Preview 2 build 3156, to support .NET Core 1.0.5 for *project.json* projects.
+ 1.1.9, which supports .NET Core 1.0.11 and 1.1.8 for *.csproj* projects.
+ 2.1.300, which supports .NET Core 2.0.0 through 2.1.0.

Any arguments passed to this builder step will be passed directly to the
`dotnet` tool. Note that when using the `dotnet restore` command you should
always specify the `--packages` option to a path within the `/workspace`, this
way the restored packages will be available to the next invocations of `dotnet`.
See the [example](examples/TestApp/cloudbuild.yaml) for how use `dotnet restore`
correctly.

We strongly advice that projects still using the `project.json` build system
contain a `global.json` file that specifies what version of the SDK is to be used.
You need make sure that you specify the right SDK so an SDK that supports `project.json`
is used to build your project.

## Building this Container Build step

To build this builder, run the following command on this directory:

```bash
gcloud builds submit --config=./cloudbuild.yaml .
```

This will build the Container Build step in your current GCP project.
