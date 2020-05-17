# curl

## Deprecation Notice

This image is deprecated.

For best support of `curl` please use one of the official `curl` images
maintained by the curlimages community on Dockerhub.

For details, visit https://hub.docker.com/r/curlimages/curl.

Note that `curlimages/curl` executes as special user `curl_user` and thus may
not echo be suitable for all purposes.

Alternatively, image `launcher.gcr.io/google/ubuntu1604` is maintained by
Google, has `curl` installed, and executes as `root`.

For details, visit https://console.cloud.google.com/launcher/details/google/ubuntu1604

## Examples

The following examples demonstrate build requests that use this builder.

### Fetch the contents of a remote URL

This `cloudbuild.yaml` fetches contents of a file by URL. For this to work the
file must be publicly readable, since no credentials are passed in the request.

```
steps:
- name: 'launcher.gcr.io/google/ubuntu1604'
  entrypoint: 'curl'
  args: ['http://www.example.com/']
```

```
steps:
- name: 'curlimages/curl'
  args: ['http://www.example.com/']
```

### Ping a remote URL

This `cloudbuild.yaml` sends a `POST` request to a URL to notify that the build
has happened, including the build's unique ID in the JSON body of the request.

```
steps:
- name: 'curlimages/curl'
  args: ['-d', '"{\"id\":\"$BUILD_ID\"}"', '-X', 'POST', 'http://www.example.com']
```
