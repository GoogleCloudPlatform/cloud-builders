# gsutil

The `gcr.io/cloud-builders/gsutil` image is maintained by the Cloud Build team,
but it may not support the most recent versions of `gsutil`. We also do not
provide historical pinned versions of `gsutil` nor support across multiple
platforms.

The Cloud SDK team maintains a `cloud-sdk` image that supports multiple tagged
versions of Cloud SDK across multiple platforms. Please visit
https://github.com/GoogleCloudPlatform/cloud-sdk-docker for details.

Suggested alternative images include:

    gcr.io/google.com/cloudsdktool/cloud-sdk
    gcr.io/google.com/cloudsdktool/cloud-sdk:alpine
    gcr.io/google.com/cloudsdktool/cloud-sdk:debian_component_based
    gcr.io/google.com/cloudsdktool/cloud-sdk:slim
    google/cloud-sdk
    google/cloud-sdk:alpine
    google/cloud-sdk:debian_component_based
    google/cloud-sdk:slim

Please note that the `gsutil` entrypoint must be specified to use these images.

When executed in the Cloud Build environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/cloud-build/docs/permissions) for the build
project.

To migrate to the Cloud SDK team's official image, make the following changes
to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/gsutil'
+ name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
+ entrypoint: 'gsutil'
```

-------

## Examples

The following examples demonstrate build requests that use `gsutil`.

For these to work, the builder service account must have permission to access
the necessary buckets and objects.

### Copy an object from Google Cloud Storage

This `cloudbuild.yaml` invokes `gsutil cp` to copy an object to the build's
workspace.

```
steps:
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
  args: ['gsutil', 'cp', 'gs://mybucket/remotefile.zip', 'localfile.zip']
```

### Copy a local file to Google Cloud Storage

This `cloudbuild.yaml` invokes `gsutil cp` to copy a local file to Google Cloud
Storage.

```
steps:
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
  args: ['gsutil', 'cp', 'localfile.zip', 'gs://mybucket/remotefile.zip']
```
