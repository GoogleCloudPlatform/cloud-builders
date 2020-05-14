# Tool builder: `gcr.io/cloud-builders/gradle`

This Cloud Build builder is deprecated and will be deleted in an upcoming
release.

Please use an [official `gradle` image](https://hub.docker.com/_/gradle/):

```yaml
steps:
- name: 'gradle'
  entrypoint: 'gradle'
  args: ['...']
```

This allows you to use any supported version of Gradle with any supported JDK
version.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
