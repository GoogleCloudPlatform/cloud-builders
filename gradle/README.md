# Tool builder: `gcr.io/cloud-builders/gradle`

This Cloud Build builder runs Gradle.

You should consider instead using an [official `gradle`
image](https://hub.docker.com/_/gradle/) and specifying the `gradle` entrypoint:

```yaml
steps:
- name: gradle:5.1.1-jdk11-slim
  entrypoint: 'gradle'
  args: ['build']
```

This allows you to use any supported version of Gradle with any supported JDK
version.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
