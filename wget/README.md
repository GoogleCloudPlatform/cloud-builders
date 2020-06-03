# wget

The `gcr.io/cloud-builders/wget` image is maintained by the Cloud Build team,
but it may not support the most recent features or versions of `wget`.

This builder simply invokes the [`wget`](https://www.gnu.org/software/wget/)
command. Arguments passed to this builder will be passed to `wget` directly.

Substantially similar functionality can be found using `curl`. While `curl`
does not offer `wget`'s ability to recursively traverse a website, `curl`
offers substantially more options for more internet protocols. `curl` is also
available in a variety of versions across multiple platforms in
community-maintained `curl` images; see the [`curl
README`](https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/curl)
for details.

Note that if `curl` is not a better option, `wget` is available in the official
community-supported [`alpine`](https://hub.docker.com/_/alpine) and
[`busybox`](https://hub.docker.com/_/busybox) images on Dockerhub, both of which
provide a variety of tagged versions.

## Examples

The following examples demonstrate build requests that use this builder.

### Fetch the contents of a remote URL

This `cloudbuild.yaml` fetches contents of a file by URL. For this to work the
file must be publicly readable, since no credentials are passed in the request.

```
steps:
- name: 'gcr.io/cloud-builders/wget'
  args: ['-O', 'localfile.zip', 'http://www.example.com/remotefile.zip']
- name: 'alpine'
  entrypoint: 'wget'
  args: ['-O', 'localfile.zip', 'http://www.example.com/remotefile.zip']
- name: 'busybox'
  entrypoint: 'wget'
  args: ['-O', 'localfile.zip', 'http://www.example.com/remotefile.zip']
- name: 'launcher.gcr.io/google/ubuntu1604'
  entrypoint: 'curl'
  args: ['-o', 'localfile.zip', 'http://www.example.com/remotefile.zip']
```

### Ping a remote URL

This `cloudbuild.yaml` sends a `POST` request to a URL to notify that the build
has happened, including the build's unique ID in the JSON body of the request.

```
steps:
- name: 'gcr.io/cloud-builders/wget'
  args: ['-q', '--post-data="{\"id\":\"$BUILD_ID\"}"', 'http://www.example.com']
- name: 'alpine'
  entrypoint: 'wget'
  args: ['-q', '--post-data="{\"id\":\"$BUILD_ID\"}"', 'http://www.example.com']
- name: 'busybox'
  entrypoint: 'wget'
  args: ['-q', '--post-data="{\"id\":\"$BUILD_ID\"}"', 'http://www.example.com']
- name: 'launcher.gcr.io/google/ubuntu1604'
  entrypoint: 'curl'
  args: ['--data-raw', '"id=$BUILD_ID"', 'http://www.example.com']
```
