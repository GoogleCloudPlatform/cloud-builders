# Tool builder: `gcr.io/cloud-builders/yarn`

This Cloud Build builder runs the `yarn` tool.

Alternative official `node` images, including multiple tagged versions
across multiple platforms are maintained by the Node.js Docker Team.

Please note that the `yarn` entrypoint must be specified when using these
images.

For further details, please visit https://hub.docker.com/_/node.

Example `cloudbuild.yaml`:

```yaml
steps:
- name: 'node:10.15.1.
  entrypoint: 'yarn'
  args: ['install']
```
