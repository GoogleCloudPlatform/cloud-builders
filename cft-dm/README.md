# cft-dm

This is a tool builder to simply invoke
[`cft`](https://github.com/GoogleCloudPlatform/deploymentmanager-samples/blob/master/community/cloud-foundation/docs/userguide.md) commands.

This tool builder invokes the Deployment Manager version of the [`toolkit`](https://cloud.google.com/foundation-toolkit/).

Arguments passed to this builder will be passed to `cft` directly, allowing
callers to run [any `cft`
command](https://github.com/GoogleCloudPlatform/deploymentmanager-samples/blob/master/community/cloud-foundation/docs/userguide.md#cli-usage).

When executed in the Cloud Build environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/cloud-build/docs/permissions) for the
project.

The latest released version of `cft` is used.

## Examples

See the examples directory