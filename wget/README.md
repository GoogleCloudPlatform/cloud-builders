# wget

This is a tool build to simply invoke the
[`wget`](https://www.gnu.org/software/wget/) command.

Arguments passed to this builder will be passed to `wget` directly.

## Examples

The following examples demonstrate build requests that use this builder.

### Fetch the contents of a remote URL

This `cloudbuild.yaml` fetches contents of a file by URL. For this to work the
file must be publicly readable, since no credentials are passed in the request.

```
steps:
- name: gcr.io/cloud-builders/wget
  args: ['-O', 'localfile.zip', 'http://www.example.com/remotefile.zip']
```

### Ping a remote URL

This `cloudbuild.yaml` sends a `POST` request to a URL to notify that the build
has happened, including the build's unique ID in the JSON body of the request.

```
steps:
- name: gcr.io/cloud-builders/wget
  args: ['-q', '--post-data="{\"id\":\"$BUILD_ID\"}"', 'http://www.example.com']
```
