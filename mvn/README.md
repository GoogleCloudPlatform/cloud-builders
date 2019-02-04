# Tool builder: `gcr.io/cloud-builders/mvn`

This Cloud Build builder runs Maven.

You should consider instead using an [official `maven`
image](https://hub.docker.com/_/maven/) and specifying the `mvn` entrypoint:

```yaml
steps:
- name: maven:3.6.0-jdk-11-slim
  entrypoint: 'mvn'
  args: ['install']
```

This allows you to use any supported version of Maven with any supported JDK
version.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
