# docker-compose

This is a tool builder to simply invoke
[`docker-compose`](https://docs.docker.com/compose/) commands.

Arguments passed to this builder will be passed to `docker-compose` directly, allowing
callers to run [any docker-compose
command](https://docs.docker.com/compose/reference/overview/).

By default, the version of docker-compose that is used by this builder is `1.11.2`.

## Examples

The following examples demonstrate build requests that use this builder:

### Build and push a container image

This `cloudbuild.yaml` invokes a `docker-compose run`, followed by a build, and then
pushes the resulting image.

```
steps:
- name: gcr.io/cloud-builders/docker-compose
  args: ['run', 'app', 'test.sh']
  id:   'test'
- name: 'gcr.io/cloud-builders/docker'
  dir: 'app'
  args: ['build', '-t', 'gcr.io/$PROJECT_ID/myimage', '.']
  id:   'build'
images: ['gcr.io/$PROJECT_ID/myimage']
```
