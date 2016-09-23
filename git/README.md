# git

This is a tool builder to simply invoke `git` commands.

Arguments passed to this builder will be passed to `git` directly.

When executed in the Container Builder environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/container-builder/docs/permissions) for the
project.
