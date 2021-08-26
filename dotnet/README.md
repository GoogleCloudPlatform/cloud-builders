# Tool builder: gcr.io/cloud-builders/dotnet

The `gcr.io/cloud-builders/dotnet` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of the dotnet
toolchain. We also do not provide historical pinned versions of dotnet tooling.

This builder is a wrapper around the
[official `mcr.microsoft.com/dotnet/core/sdk:2.1` image](https://hub.docker.com/r/microsoft/dotnet/)
that differs in that it specifies `dotnet` as an entrypoint. It is functionally
equivalent to:

```yaml
steps:
- name: 'mcr.microsoft.com/dotnet/core/sdk:2.1'
  entrypoint: 'dotnet'
```

For alternative official `dotnet` builer images, including additional dotnet
tooling suitable for use in Cloud Build, please visit
https://hub.docker.com/_/microsoft-dotnet).

To migrate to the official `dotnet` image, make the following changes to your
`cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/dotnet'
+ name: 'microsoft/dotnet:sdk'
+ entrypoint: 'dotnet'
```
