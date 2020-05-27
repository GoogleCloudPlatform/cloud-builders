# Tool builder: `gcr.io/cloud-builders/mvn`

The `gcr.io/cloud-builders/mvn` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of Maven. We also do
not provide historical pinned versions of Maven.

A supported `maven` image, including multiple tagged versions across multiple
versions of Java and multiple platforms, is maintained by the Maven team at
https://hub.docker.com/_/maven.

To migrate to the Maven team's official `maven` image, make the following
changes to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/maven'
+ name: 'maven'
+ entrypoint: 'mvn'
```

## Example:

```yaml
steps:
- name: 'maven'
  entrypoint: 'mvn'
  args: ['install']
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
