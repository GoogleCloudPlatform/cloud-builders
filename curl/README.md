# curl

## Altnernative Official Images

This image is a simple wrapper on top of `launcher.gcr.io/google/ubuntu1604`
that specifies `curl` as the `entrypoint`. For a community-supported version of
`curl`, `launcher.gcr.io/google/ubuntu1604` can be used directly with Cloud Build.
For details, visit https://console.cloud.google.com/launcher/details/google/ubuntu1604.

Also note that the community-supported
[`curlimages/curl`](https://hub.docker.com/r/curlimages/curl) is compatible
with Cloud Build and available in numerous tagged versions for multiple
platforms, but it runs as user `curl_user` and thus may not be suitable for all
purposes. For details, visit https://hub.docker.com/r/curlimages/curl.


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
- name: 'launcher.gcr.io/google/ubuntu1604'
  entrypoint: 'curl'
  args: ['-d', '"{\"id\":\"$BUILD_ID\"}"', '-X', 'POST', 'http://www.example.com']
```

```
steps:
- name: 'curlimages/curl'
  args: ['-d', '"{\"id\":\"$BUILD_ID\"}"', '-X', 'POST', 'http://www.example.com']
```
