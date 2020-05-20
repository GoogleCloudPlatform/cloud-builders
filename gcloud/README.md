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

Note: official `cloud-sdk` images, including multiple tagged versions across
multiple platforms, can be found at
https://github.com/GoogleCloudPlatform/cloud-sdk-docker.

Suggested alternative images include:

    gcr.io/google.com/cloudsdktool/cloud-sdk
    gcr.io/google.com/cloudsdktool/cloud-sdk:slim
    gcr.io/google.com/cloudsdktool/cloud-sdk:alpine

Please note that the `gcloud` entrypoint must be specified to use these images.

-------

## Examples

The following examples demonstrate build requests that use `gcloud`.

### Clone a Cloud Source Repository using `gcloud`

This `cloudbuild.yaml` invokes `gcloud source repos clone` to clone the
`default` Cloud Source Repository.

```
steps:
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
  args: ['gcloud', 'source', 'repos', 'clone', 'default']
```


## `gcloud` vs `gcloud-slim`

There are two variants of the `gcloud` builder in this repository:

* `gcloud` installs all optional gcloud components, and is much larger.
* `gcloud-slim` installs only the `gcloud` CLI and no components, and is
  smaller.
