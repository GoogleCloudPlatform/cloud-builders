# gsutil

This is a tool builder to simply invoke
[`gsutil`](https://cloud.google.com/storage/docs/gsutil) commands.

Arguments passed to this builder will be passed to `gsutil` directly.

When executed in the Container Builder environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/container-builder/docs/permissions) for the
project.

The latest available version of `gsutil` is used.
