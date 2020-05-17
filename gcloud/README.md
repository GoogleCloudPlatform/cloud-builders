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

Note: This image is deprecated in favor of
[`gcr.io/google.com/cloudsdktool/cloud-sdk`](https://github.com/GoogleCloudPlatform/cloud-sdk-docker). Users can switch
to using the `gcr.io/google.com/cloudsdktool/cloud-sdk` image today for a maintained and up-to-date
`gcloud` builder. The deprecation is tracked in
[issue638](https://github.com/GoogleCloudPlatform/cloud-builders/issues/638).

If your testing with
`gcr.io/google.com/cloudsdktool/cloud-sdk` reveals incompatibilities, please post a comment in
[issue638](https://github.com/GoogleCloudPlatform/cloud-builders/issues/638).

-------

## Examples

The following examples demonstrate build request that use `gcloud`.

### Clone a Cloud Source Repository using `gcloud`

This `cloudbuild.yaml` invokes `gcloud source repos clone` to clone the
`default` Cloud Source Repository.

```
steps:
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
  args: ['gcloud', 'source', 'repos', 'clone', 'default']
```


## `gcloud` vs `gcloud-slim`

There are two variants of the `gcloud` builder:

* `gcloud` installs all optional gcloud components, and is much larger.
* `gcloud-slim` installs only the `gcloud` CLI and no components, and is
  smaller.
