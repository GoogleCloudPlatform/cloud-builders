# Tool builder: `gcr.io/cloud-builders/yarn`

This Cloud Build builder runs the `yarn` tool.

## Deprecation Notice

This builder is deprecated in favor of the supported
[official `node` images](https://hub.docker.com/_/node/).

Example `cloudbuild.yaml`:

```yaml
steps:
- name: node:10.15.1
  entrypoint: yarn
  args: ['install']
```

This builder will be deleted in an upcoming release.
