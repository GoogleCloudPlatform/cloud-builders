# Tool builder: `gcr.io/cloud-builders/javac`

The `gcr.io/cloud-builders/javac` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of `javac`. We also
do not provide historical pinned versions of Java.

A supported `javac` image, including multiple tagged versions across multiple
versions of Java and multiple platforms, is maintained by the OpenJDK team at
https://hub.docker.com/_/openjdk. These images also support additional Java
tooling.

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
