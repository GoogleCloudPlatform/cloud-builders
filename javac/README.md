# Tool builder: `gcr.io/cloud-builders/javac`

The `gcr.io/cloud-builders/javac` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of `javac`. We also
do not provide historical pinned versions of Java.

The OpenJDK team provides `openjdk` images that provide additional Java tooling
and support multiple tagged versions across multiple versions of Java and
multiple platforms. Please visit https://hub.docker.com/_/openjdk for details.

To migrate to the OpenJDK team's official image, make the following changes to
your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/javac'
+ name: 'openjdk'
+ entrypoint: 'javac'
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
