# gcloud

This is a tool builder to simply invoke
[`gcloud`](https://cloud.google.com/sdk/gcloud/) commands.

Arguments passed to this builder will be passed to `gcloud` directly.

When executed in the Container Builder environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/container-builder/docs/permissions) for the
project.

The latest available version of `gcloud` is used.
