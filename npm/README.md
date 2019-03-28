# Tool builder: `gcr.io/cloud-builders/npm`

This Cloud Build builder runs the `npm` tool.

You should consider instead using an [official `node`
image](https://hub.docker.com/_/node/) and specifying the `npm` entrypoint:

```yaml
steps:
- name: node:10.15.1
  entrypoint: npm
  args: ['install']
```

This allows you to use any supported version of NPM.

## Building this builder

To build this builder, run the following command in this directory.

```bash
$ gcloud builds submit
```