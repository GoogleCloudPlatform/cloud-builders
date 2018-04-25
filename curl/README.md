# curl

This is a tool build to simply invoke the
[`curl`](https://www.gnu.org/software/curl/) command.

Arguments passed to this builder will be passed to `curl` directly.

## Examples

The following examples demonstrate build requests that use this builder.

### Fetch the contents of a remote URL

This `cloudbuild.yaml` fetches contents of a file by URL. For this to work the
file must be publicly readable, since no credentials are passed in the request.

```
steps:
- name: gcr.io/cloud-builders/curl
  args: ['http://www.example.com/remotefile.zip', '--output', 'localfile.zip']
```

### Ping a remote URL

This `cloudbuild.yaml` sends a `POST` request to a URL to notify that the build
has happened, including the build's unique ID in the JSON body of the request.

```
steps:
- name: gcr.io/cloud-builders/curl
  args: ['-d', '"{\"id\":\"$BUILD_ID\"}"', '-X', 'POST', 'http://www.example.com']
```
