# Tool builder: `gcr.io/cloud-builders/mvn`

This Cloud Build builder runs Maven. It also includes a number of dependencies
that are precached within the image.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit . --config=cloudbuild.yaml
