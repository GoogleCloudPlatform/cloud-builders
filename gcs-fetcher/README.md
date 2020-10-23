# GCS Fetcher

** Warning: This builder is experimental and is very likely to change in
breaking ways at this time. **

This tool fetches objects from Google Cloud Storage, either in the form of a
.zip archive, or based on the contents of a source manifest file.

## Source Manifests

A source manifest is a JSON object in Cloud Storage listing *other* objects in
Cloud Storage that should be fetched. The format of the manifest is a mapping of
destination file path to the location in Cloud Storage where the file's contents
can be found.

The following is an example source manifest:

```json
{
  "Dockerfile": {
    "sourceUrl": "gs://my-bucket/abcdef",
    "sha1sum": "<sha-1 digest>"
  },
  "path/to/main.go": {
    "sourceUrl": "gs://my-bucket/ghijk",
    "sha1sum": "<sha-1 digest>"
  }
}
```

To process the above manifest, the GCS Fetcher tool processes each element:

1. Fetch the object located at `sourceUrl`
1. Verify the object's SHA-1 matches the expected digest
1. Write the file contents to the path indicated by the object key

So in the above example, the tool fetches `gs://my-bucket/abcdef`, verifies its
SHA-1 digest, and places the file in the working directory named as
`Dockerfile`. It then fetches `gs://my-bucket/ghijk`, verifies its SHA-1 digest,
and places the file in the working directory at `path/to/main.go`.

### Why Source Manifests?

The main benefit to source manifests are in enabling incremental upload of
sources from a client.

If a user uploads the two files in the example above, then edits the contents of
their `Dockerfile`, only that new file has to be uploaded to Cloud Storage, and
the manifest for the next build can reuse the entry for `path/to/main.go`,
only writing a new entry for `Dockerfile` to specify the new file's location in
Cloud Storage.

In the case of a file move, where file contents don't change, no new files have
to be uploaded. The new manifest simply changes the key in the object describing
the path to place the file fetched from Cloud Storage.

## Full Example

To fetch source described in a source manifest, add the following line to your
build config:

```yaml
steps:
- name: 'gcr.io/cloud-builders/gcs-fetcher'
  args:
  - '--type=Manifest'
  - '--location=gs://${PROJECT_ID}_cloudbuild/manifest-foo.json'
```


It may also be useful to _produce_ and upload source manifests describing some
source, which you can do with `gcr.io/cloud-builders/gcs-uploader`:

```yaml
steps:
- name: 'gcr.io/cloud-builders/gcs-uploader'
  args: ['--location=${PROJECT_ID}_cloudbuild/manifest-${BUILD_ID}.json']
```

This will upload the contents of the workspace directory, ignoring objects that
are already present in Cloud Storage, and upload a manifest JSON object named
`manifest-${BUILD_ID}.json` to the same Cloud Storage bucket.

`gcs-uploader` will not delete remote objects that are not present locally.

### Caching resources

`gcs-fetcher` and `gcs-uploader` can be used together to provide simple
cross-build caching functionality, by optimistically fetching files that are
expensive to generate at the beginning of the build, and by uploading those
files at the end of the build.

In order to benefit from this you would need to define a well-known reusable
manifest file location.

```yaml
steps:
# Attempt to fetch whatever files are available.
- name: 'gcr.io/cloud-builders/gcs-fetcher'
  args:
  - '--type=Manifest'
  - '--location=gs://${PROJECT_ID}_cloudbuild_cache/manifest-foo.json'

# Generate new files, ignoring those that already exist.
# - name: 'generate-new-files'
#   ...

# Upload all files; only new and changed files will be uploaded to
# Cloud Storage.
- name: 'gcr.io/cloud-builders/gcs-uploader'
  - '--bucket=${PROJECT_ID}_cloudbuild_cache'
  - '--manifest_file=manifest-foo.json'
```

**Tips and Caveats**

1. You can use [Cloud Storage object lifecycle
   management](https://cloud.google.com/storage/docs/lifecycle) to automatically
   delete objects after a certain amount of time, which may reduce storage costs
   but negatively impact cache hit rates.
1. Two ongoing builds that fetch from and upload to the same manifest file may
   interact poorly with each other, leading to confusing bugs. You may want to
   include `${BRANCH_NAME}` or some other unique value in the manifest file
   location to avoid this.
1. Even with incremental upload, you may find that generating the files is just
   as fast or faster than fetching from Cloud Storage. Caching is not a magic
   bullet, and can add more complexity than it removes.

## Outstanding TODOs:

- [ ] .tar.gz support, depending on object name extension
- [ ] Unit tests to cover generation parsing/formatting
- [ ] Actually verify SHA-1
