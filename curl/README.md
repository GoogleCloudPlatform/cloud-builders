# curl

The `gcr.io/cloud-builders/curl` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of `curl`. We also do
not provide tagged versions or support for multiple OS platforms.

The `curl` community maintains a `curl` image that supports multiple tagged
versions at [`curlimages/curl`](https://hub.docker.com/r/curlimages/curl). While
this image is compatible with the hosted Cloud Build service, it runs as user
`curl_user` and thus may not be suitable for all purposes. For details, visit
https://hub.docker.com/r/curlimages/curl.

This `gcr.io/cloud-builders/curl` image is a simple wrapper on top of
`gcr.io/gcp-runtimes/ubuntu_18_0_4` that specifies `curl` as the `entrypoint`.
As a Google-provided image, `gcr.io/gcp-runtimes/ubuntu_18_0_4` can be used
directly with Cloud Build.  For details, visit
https://console.cloud.google.com/marketplace/product/google/ubuntu1804. Using this
image directly will mean that you are always using the latest patched version.

To migrate to the GCP launcher image, make the following changes
to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/curl'
+ name: 'gcr.io/gcp-runtimes/ubuntu_18_0_4'
+ entrypoint: 'curl'
```

## Examples

The following examples demonstrate build requests that use this builder.

### Fetch the contents of a remote URL

This `cloudbuild.yaml` fetches contents of a file by URL. For this to work the
file must be publicly readable, since no credentials are passed in the request.

```
steps:
- name: 'gcr.io/gcp-runtimes/ubuntu_18_0_4'
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
- name: 'gcr.io/gcp-runtimes/ubuntu_18_0_4'
  entrypoint: 'curl'
  args: ['-d', '"{\"id\":\"$BUILD_ID\"}"', '-X', 'POST', 'http://www.example.com']
```

```
steps:
- name: 'curlimages/curl'
  args: ['-d', '"{\"id\":\"$BUILD_ID\"}"', '-X', 'POST', 'http://www.example.com']
```
