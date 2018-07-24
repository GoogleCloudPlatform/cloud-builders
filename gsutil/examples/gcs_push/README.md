# GCS Push Example

This example demonstrates a build that:

1.  Builds a "Hello, world" Go binary named `hello`;
2.  Copies the binary to a GCS bucket.

To run this example, make sure you have created a GCS Bucket named `$PROJECT_ID`
and run:

```
gcloud builds submit --config=cloudbuild.yaml .
```

The `hello` binary will be found in your GCS Bucket.
