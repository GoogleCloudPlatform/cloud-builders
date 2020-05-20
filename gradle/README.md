# Tool builder: `gcr.io/cloud-builders/gradle`

Alternative official \`gradle\` images, including multiple taged version across
multiple jdk versions and platforms, can be found at
https://hub.docker.com/_/gradle.

```yaml
steps:
- name: 'gradle'
  entrypoint: 'gradle'
  args: ['...']
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
