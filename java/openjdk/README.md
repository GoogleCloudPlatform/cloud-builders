# Tool builder: `gcr.io/cloud-builders/java/openjdk`

This build step contains openjdk and docker. It is used as a base image for other java tool
builders.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud container builds submit . --config=cloudbuild.yaml
