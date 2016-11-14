# gcloud

This is a tool builder to simply invoke
[`gcloud`](https://cloud.google.com/sdk/gcloud/) commands.

Arguments passed to this builder will be passed to `gcloud` directly, allowing
callers to run [any `gcloud`
command](https://cloud.google.com/sdk/gcloud/reference/).

When executed in the Container Builder environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/container-builder/docs/permissions) for the
project.

The latest released version of `gcloud` is used.

## Examples

The following examples demonstrate build request that use this builder:

### Clone a Cloud Source Repository using `gcloud`

This `cloudbuild.yaml` invokes `gcloud source repos clone` to clone the
`default` Cloud Source Repository.

```
steps:
- name: gcr.io/cloud-builders/gcloud
  args: ['source', 'repos', 'clone', 'default']
```
