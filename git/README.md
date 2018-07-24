# git

This is a tool builder to simply invoke `git` commands.

Arguments passed to this builder will be passed to `git` directly, allowing
callers to run [any Git command](https://git-scm.com/docs).

When executed in the Cloud Build environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/cloud-build/docs/permissions) for the
project.

## Examples

The following examples demonstrate build requests that use this builder.

For these to work, the remote repository must either be public, or the builder
service account must have permission to access it.

### Clone a Git repository

This `cloudbuild.yaml` clones a remote Git repository to the build's workspace.

```
steps:
- name: gcr.io/cloud-builders/git
  args: ['clone', 'https://github.com/GoogleCloudPlatform/cloud-builders']
```

### Push changes to a remote Git repository

This `cloudbuild.yaml` pushes local changes to a remote Git repository.

```
steps:
- name: gcr.io/cloud-builders/git
  args: ['push', 'https://source.developers.google.com/p/$PROJECT_ID/r/myrepo', 'master']
```
