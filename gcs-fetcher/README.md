# GCS Fetcher

** This builder is experimental and currently very likely to change in breaking
ways at this time. **

This tool fetches objects from Google Cloud Storage, either in the form of a
.zip archive, or based on the contents of a source manifest file.

## Source Manifests

A source manifest is a JSON object in Cloud Storage listing *other* objects in
Cloud Storage which should be fetched. The format of the file is a mapping of
destination file path to the location in GCS where the file's contents can be
found.

For example:

```
{
  "Dockerfile": {
    "sourceUrl": "gs://my-bucket/abcdef",
    "sha1sum": "<sha-1 digest>"
  },
  "path/to/main.go": {
    "sourceUrl": "gs://my-bucket/ghijk",
    "sha1sum": "<sha-1 digest>"
  },
  ...
}
```

To process this manifest, the tool each element:

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

In the case of a file move, where file contents themselves don't change, no new
files have to be uploaded, the new manifest simply changes the key in the object
describing the path to place the file fetched from Cloud Storage.

## Outstanding TODOs:

- [ ] .tar.gz support, depending on object name extension
- [ ] Simple manifest upload tool, for manual testing
- [ ] Unit tests to cover generation parsing/formatting
- [ ] Actually verify SHA-1
