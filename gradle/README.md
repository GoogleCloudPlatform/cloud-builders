# Tool builder: `gcr.io/cloud-builders/gradle`

This Cloud Build builder runs the Gradle build tool.

You should consider instead using an [official `gradle`
image](https://hub.docker.com/_/gradle/) and specifying the `gradle` entrypoint:
```yaml
steps:
- name: gradle:6.0.1-jdk11
  entrypoint: 'gradle'
  args: ['build']
```
This allows you to use any supported version of Gradle with any supported JDK
version.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit . --config=cloudbuild.yaml
