# Tool builder: `gcr.io/cloud-builders/yarn`

This Cloud Build builder runs the `yarn` tool.

You might also consider using an [official `node` image](https://hub.docker.com/_/node/) and specifying the `yarn` entrypoint:

```yaml
steps:
- name: node:10.10.0
  entrypoint: yarn
  args: ['install']
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit . --config=cloudbuild.yaml
