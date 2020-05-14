# Tool builder: `gcr.io/cloud-builders/mvn`

This Cloud Build builder is deprecated and will be deleted in an upcoming
release.

Please use an [official `maven` image](https://hub.docker.com/_/maven/):

```yaml
steps:
- name: 'maven'
  entrypoint: 'mvn'
  args: ['install']
```

This allows you to use any supported version of Maven with any supported JDK
version.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
