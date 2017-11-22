# Tool builder: `gcr.io/cloud-builders/scala/sbt`

This Container Builder build step runs the [sbt](http://www.scala-sbt.org/).

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud container builds submit . --config=cloudbuild.yaml
