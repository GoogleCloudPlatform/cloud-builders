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

Note: This image is deprecated in favor of
[`google/cloud-sdk`](https://hub.docker.com/r/google/cloud-sdk/). Users can switch
to using the `google/cloud-sdk` image today for a maintained and up-to-date
`gsutil` builder. The deprecation is tracked in
[issue638](https://github.com/GoogleCloudPlatform/cloud-builders/issues/638).

If your testing with
`google/cloud-sdk` reveals incompatibilities, please post a comment in
[issue638](https://github.com/GoogleCloudPlatform/cloud-builders/issues/638).

## Examples

The following examples demonstrate build requests that use `gsutil`.

For these to work, the builder service account must have permission to access
the necessary buckets and objects.

### Copy an object from Google Cloud Storage

This `cloudbuild.yaml` invokes `gsutil cp` to copy an object to the build's
workspace.

```
steps:
- name: 'google/cloud-sdk'
  args: ['gsutil', 'cp', 'gs://mybucket/remotefile.zip', 'localfile.zip']
```

### Copy a local file to Google Cloud Storage

This `cloudbuild.yaml` invokes `gsutil cp` to copy a local file to Google Cloud
Storage.

```
steps:
- name: 'google/cloud-sdk'
  args: ['gsutil', 'cp', 'localfile.zip', 'gs://mybucket/remotefile.zip']
```
