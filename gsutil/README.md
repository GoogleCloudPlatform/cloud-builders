# gsutil

This is a tool builder to simply invoke
[`gsutil`](https://cloud.google.com/storage/docs/gsutil) commands.

Arguments passed to this builder will be passed to `gsutil` directly, allowing
callers to run [any `gsutil`
command](https://cloud.google.com/storage/docs/gsutil).

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

Please note that the `gsutil` entrypoint must be specified to use these images.

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
