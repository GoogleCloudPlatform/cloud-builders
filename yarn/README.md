# Tool builder: `gcr.io/cloud-builders/yarn`

This Cloud Build builder runs the `yarn` tool.

You might also consider using the [official `node` image] and specifying the
`yarn` entrypoint:

```yaml
steps:
- name: node
  entrypoint: yarn
  args: ['install']
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit . --config=cloudbuild.yaml
