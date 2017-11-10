# Tool builder: `gcr.io/cloud-builders/mvn`

This Container Builder build step runs Maven. It also includes a number of
dependencies that are precached within the image.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud container builds submit . --config=cloudbuild.yaml
