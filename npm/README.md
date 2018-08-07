# Tool builder: `gcr.io/cloud-builders/npm`

This Cloud Build builder runs the `npm` tool.

You might also consider using the [official `node` image] and specifying the
`npm` entrypoint:

```yaml
steps:
- name: node
  entrypoint: npm
  args: ['install']
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit . --config=cloudbuild.yaml
