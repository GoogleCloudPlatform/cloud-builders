# Tool builder: `gcr.io/cloud-builders/npm`

This Cloud Build builder runs the `npm` tool.

## Deprecation Notice

This builder is deprecated in favor of the supported
[official `node` images](https://hub.docker.com/_/node/).

Example `cloudbuild.yaml`:

```yaml
steps:
- name: node
  entrypoint: npm
  args: ['install']
```

This builder will be deleted in an upcoming release.
