# Tool builder: `gcr.io/cloud-builders/gcloud`

This Cloud Build builder runs the
[`gcloud`](https://cloud.google.com/sdk/gcloud/) tool to interact with Google
Cloud Platform resources.

When executed in the Cloud Build environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/cloud-build/docs/permissions) for the
project. The latest released version of `gcloud` is used.

You should consider instead using an [official `google/cloud-sdk`
image](https://hub.docker.com/r/google/cloud-sdk) and specifying the `gcloud` entrypoint:

```yaml
steps:
- name: google/cloud-sdk:230.0.0
  entrypoint: 'gcloud'
  args: ['version']
```

This allows you to use any supported version of `gcloud`.

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

* `gcloud` has all optional gcloud components installed, and is much larger.
* `gcloud-slim` installs only the `gcloud` CLI and no components, and is
  smaller.

Both images are cached on Cloud Build VMs, so the size of the image should not
matter in most cases when running in that environment. However, in other
environments where images are not cached, you may find that a smaller builder
image is faster to pull, and might be preferrable to the larger "kitchen sink"
`gcloud` builder image.
