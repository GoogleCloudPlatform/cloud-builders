# gsutil

This is a tool builder to simply invoke
[`gsutil`](https://cloud.google.com/storage/docs/gsutil) commands.

Arguments passed to this builder will be passed to `gsutil` directly, allowing
callers to run [any `gsutil`
command](https://cloud.google.com/storage/docs/gsutil).

When executed in the Container Builder environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/container-builder/docs/permissions) for the
project.

The latest available version of `gsutil` is used.

## Examples

The following examples demonstrate build request that use this builder:

For these to work, the builder service account must have permission to access
the necessary buckets and objects.

### Copy an object from Google Cloud Storage

This `cloudbuild.yaml` invokes `gsutil cp` to copy an object to the build's
workspace.

```
steps:
- name: gcr.io/cloud-builders/gsutil
  args: ['cp', 'gs://mybucket/remotefile.zip', 'localfile.zip']
```

### Copy a local file to Google Cloud Storage

This `cloudbuild.yaml` invokes `gsutil cp` to copy a local file to Google Cloud
Storage.

```
steps:
- name: gcr.io/cloud-builders/gsutil
  args: ['cp', 'localfile.zip', 'gs://mybucket/remotefile.zip']
```
