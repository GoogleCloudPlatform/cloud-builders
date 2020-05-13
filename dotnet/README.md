# Tool builder: gcr.io/cloud-builders/dotnet

This builder is deprecated and will be deleted in an upcoming release.

Please use one of the [official `microsoft/dotnet:sdk`
images](https://hub.docker.com/r/microsoft/dotnet/) and specifying the `dotnet`
entrypoint:

```yaml
steps:
- name: 'microsoft/dotnet:sdk'
  entrypoint: 'dotnet'
  args: ['build']
```
