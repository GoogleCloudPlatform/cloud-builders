# Tool builder: `gcr.io/cloud-builders/yarn`

The `gcr.io/cloud-builders/yarn` image is maintained by the Cloud Build team, but
it may not support the most recent features or versions of `yarn`. We also do not
provide historical pinned versions of `yarn`.

A supported `yarn` image, including multiple tagged versions, is maintained by
the Node team at https://hub.docker.com/_/node. This image also provides
additional Node tooling.

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
