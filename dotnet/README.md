# Tool builder: gcr.io/cloud-builders/dotnet

This Cloud Build builder runs the `dotnet` SDK CLI.

You might also consider using an [official `microsoft/dotnet:sdk`
image](https://hub.docker.com/microsoft/dotnet/) and specifying the `dotnet`
entrypoint:

```yaml
steps:
- name: microsoft/dotnet:sdk
  entrypoint: dotnet
  args: ['build']
```

The `gcr.io/cloud-builders/dotnet` image provides the latest released version of
the dotnet CLI, exactly what's tagged as `microsoft/dotnet:sdk`. To use another
version of the tool, specify a tagged image with that version, e.g.:

```yaml
steps:
- name: microsoft/dotnet:1.1-sdk
  entrypoint: dotnet
  args: ['build']
```

Any arguments passed to this builder will be passed directly to the `dotnet`
tool. Note that when using the `dotnet restore` command you should always
specify the `--packages` option to a path within the `/workspace`, this way the
restored packages will be available to the next invocations of `dotnet`.

## Building this Container Build step

To build this builder, run the following command on this directory:

```bash
gcloud builds submit --config=./cloudbuild.yaml
```
