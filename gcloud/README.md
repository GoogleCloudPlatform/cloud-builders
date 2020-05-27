# gcloud

The `gcr.io/cloud-builders/gcloud` image is maintained by the Cloud Build team,
but it may not support the most recent versions of `gcloud`. We also do not
provide historical pinned versions of `gcloud` not support across multiple
platforms.

A supported `cloud-sdk` image, including multiple tagged versions across several
platforms, is maintained by the Cloud SDK team at
https://github.com/GoogleCloudPlatform/cloud-sdk-docker.

Suggested alternative images include:

    gcr.io/google.com/cloudsdktool/cloud-sdk
    gcr.io/google.com/cloudsdktool/cloud-sdk:slim
    gcr.io/google.com/cloudsdktool/cloud-sdk:alpine
	google/cloud-sdk
	google/cloud-sdk:slim
	google/cloud-sdk:alpine

Please note that the `gcloud` entrypoint must be specified to use these images.

When executed in the Cloud Build environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/cloud-build/docs/permissions) for the build
project.

To migrate to the Cloud SDK team's official image, make the following changes
to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/gcloud'
+ name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
+ entrypoint: 'gcloud'
```

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
