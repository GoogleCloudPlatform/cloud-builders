# Tool builder: gcr.io/cloud-builders/dotnet

This builder is a wrapper around the [official `microsoft/dotnet:sdk`
images](https://hub.docker.com/r/microsoft/dotnet/) that specifies `dotnet` as
an entrypoint. It is functionally equivalent to:

```yaml
steps:
- name: 'microsoft/dotnet:sdk'
  entrypoint: 'dotnet'
  args: ['build']
```

For an alternative official `dotnet` builer images, including multiple tagged
versions across a variety of platforms, please visit
https://hub.docker.com/r/microsoft/dotnet/).
