# Tool builder: `gcr.io/cloud-builders/mvn`

The `gcr.io/cloud-builders/mvn` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of Maven. We also do
not provide historical pinned versions of Maven.

The Maven team provides a `maven` image that supports multiple tagged versions
of Maven across multiple versions of Java and multiple platforms. Please visit
https://hub.docker.com/_/maven for details.

To migrate to the Maven team's official `maven` image, make the following
changes to your `cloudbuild.yaml`:

```diff
- name: 'gcr.io/cloud-builders/mvn'
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

## Python, Cloud SDK, and App Engine use

The community-supported `maven` image no longer ships with `python`; this breaks
the [App Engine `maven`
plugin](https://cloud.google.com/appengine/docs/standard/java/using-maven).

This builder supports
[`Dockerfile.appengine`](Dockerfile.appengine) that is compatible with the
App Engine-compatible image and is hosted at
`gcr.io/cloud-builders/mvn:appengine`.

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
