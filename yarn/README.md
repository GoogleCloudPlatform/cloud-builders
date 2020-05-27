# Tool builder: `gcr.io/cloud-builders/yarn`

The `gcr.io/cloud-builders/yarn` image is maintained by the Cloud Build team, but
it may not support the most recent features or versions of `yarn`. We also do not
provide historical pinned versions of `yarn`.

The Node team provides `node` images that support multiple tagged versions of
`yarn` and additional Node tooling. Please visit https://hub.docker.com/_/node
for details.

To migrate to the Node team's official Node image, make the following changes
to your `cloudbuild.yaml`:

```
- name: 'gcr.io/cloud-builders/yarn'
+ name: 'node'
+ entrypoint: 'yarn'
```


## Example:

`cloudbuild.yaml`:

```yaml
steps:
- name: 'node'
  entrypoint: 'yarn'
  args: ['install']
```
