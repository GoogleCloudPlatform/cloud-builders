# Tool builder: `gcr.io/cloud-builders/npm`

The `gcr.io/cloud-builders/npm` image is maintained by the Cloud Build team, but
it may not support the most recent features or versions of `npm`. We also do not
provide historical pinned versions of `npm`.

The Node team provides `node` images that support multiple tagged versions of
`npm` and additional Node tooling. Please visit https://hub.docker.com/_/node
for details.

To migrate to the Node team's official Node image, make the following changes
to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/npm'
+ name: 'node'
+ entrypoint: 'npm'
```


## Example:

`cloudbuild.yaml`:

```yaml
steps:
- name: 'node'
  entrypoint: 'npm'
  args: ['install']
```
