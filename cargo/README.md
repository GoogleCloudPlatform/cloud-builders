# Tool builder: `gcr.io/cloud-builders/cargo`

This Container Builder build step runs the `cargo` tool.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud container builds submit . --config=cloudbuild.yaml
