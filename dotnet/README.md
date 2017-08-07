# Tool builder: gcr.io/cloud-builders/dotnet
This directory contains the description for the `dotnet` Container Build step. This step will run the .NET Core `dotnet` command line. This builder has two SDKs installed:
+ 1.0.3 with SDK Preview 2 build 3156, to support .NET Core 1.0.5 with project.json solutions.
+ 1.0.4, which supporst .csproj projects for .NET Core 1.0.5 and 1.1.2.

Any arguments passed to this builder step will be passed directly to the `dotnet` tool. Note that when using the `dotnet restore` command you should always specify the `--packages` option to a path within the `/workspace`, this way the restored packages will be available to the next invocations of `dotnet`. See the [example](examples/TestApp/cloudbuild.yaml) for how use `dotnet resotre` correctly.

## Building this Container Build step
To build this builder, run the following command on this directory:
```bash
gcloud container builds submit --config=./cloudbuild.yaml .
```

This will build the Container Build step in your current GCP project.
