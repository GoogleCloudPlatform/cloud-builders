# Tool builder: `gcr.io/cloud-builders/gradle`

The `gcr.io/cloud-builders/gradle` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of Gradle. We also do
not provide historical pinned versions of Gradle.

The Gradle team provides a `gradle` image that supports multiple tagged versions
across multiple versions of Java and multiple platforms. Please visit
https://hub.docker.com/_/gradle for details.

To migrate to the Gradle team's official `gradle` image, make the following
changes to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/gradle'
+ name: 'gradle'
+ entrypoint: 'gradle'
```

## Example:

```yaml
steps:
- name: 'gradle'
  entrypoint: 'gradle'
  args: ['...']
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
