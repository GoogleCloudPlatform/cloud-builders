# gcloud

This is a tool builder to simply invoke
[`gcloud`](https://cloud.google.com/sdk/gcloud/) commands.

Arguments passed to this builder will be passed to `gcloud` directly, allowing
callers to run [any `gcloud`
command](https://cloud.google.com/sdk/gcloud/reference/).

When executed in the Cloud Build environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/cloud-build/docs/permissions) for the
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


## `gcloud` vs `gcloud-slim`

There are two variants of the `gcloud` builder:

* `gcloud` installs all optional gcloud components, and is much larger.
* `gcloud-slim` installs only the `gcloud` CLI and no components, and is
  smaller.

Both images are cached on Cloud Build VMs, so the size of the image should not
matter in most cases when running in that environment. However, in other
environments where images are not cached, you may find that a smaller builder
image is faster to pull, and might be preferrable to the larger "kitchen sink"
`gcloud` builder image.
