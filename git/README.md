# git

This is a tool builder to simply invoke `git` commands.

Arguments passed to this builder will be passed to `git` directly, allowing
callers to run [any Git command](https://git-scm.com/docs).

When executed in the Cloud Build environment, commands are executed with
credentials of the [builder service
account](https://cloud.google.com/cloud-build/docs/permissions) for the
project. This enables access to [Cloud Source
Repositories](https://cloud.google.com/source-repositories) when used in builds
on the hosted Cloud Build service.

The `gcr.io/cloud-builders/git` image is maintained by the Cloud Build team, but
it may not support the most recent versions of `git`. Please note that the
`cloud-sdk` images supported by the Cloud SDK team similarly include a `git`
installation configured to use [Application Default
Credentials](https://cloud.google.com/docs/authentication/production) when
running in the hosted Cloud Build service; these images also provide support for
a wider combination of platform and Cloud SDK version variations and thus may be
more suitable for some use cases when interacting with Cloud Source
Repositories.

Suggested alternative images include:

    gcr.io/google.com/cloudsdktool/cloud-sdk
    gcr.io/google.com/cloudsdktool/cloud-sdk:alpine
    gcr.io/google.com/cloudsdktool/cloud-sdk:debian_component_based
    gcr.io/google.com/cloudsdktool/cloud-sdk:slim
    google/cloud-sdk
    google/cloud-sdk:alpine
    google/cloud-sdk:debian_component_based
    google/cloud-sdk:slim

These images are automatically configured to use Application Default Credentials
in versions `295.0.0` and later.

Please note that the `git` entrypoint must be specified to use these images.

## Examples

The following examples demonstrate build requests that use this builder.

For these to work, either the remote repository must be public or the builder
service account must have permission to access it.

### Clone a Git repository

This `cloudbuild.yaml` demonstrates cloning a Git repository to the build's
workspace.

```
steps:
# This build step clones a remote public repository from Github; no
# credentials are required.
- name: 'gcr.io/cloud-builders/git'
  args: ['clone', 'https://github.com/GoogleCloudPlatform/cloud-builders']
# This step uses Application Default Credentials to clone a Cloud Source
# Repository.
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk:alpine'
  entrypoint: 'git'
  args: ['clone', 'https://source.developers.google.com/p/$PROJECT_ID/r/$REPO']
# This step is functionally equivalent to the prior step, but invokes `gcloud`
# instead of using `git` directly.
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
  entrypoint: 'gcloud'
  args: ['source', 'repos', 'clone', '$REPO']
```

### Push changes to a remote Git repository

This `cloudbuild.yaml` demonstrates two functionally equivalent ways to push
local changes to a remote Git repository authenticated with Application Default
Credentials.

```
steps:
- name: 'gcr.io/cloud-builders/git'
  args: ['push', 'https://source.developers.google.com/p/$PROJECT_ID/r/$REPO', 'master']
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk:alpine'
  entrypoint: 'git'
  args: ['push', 'https://source.developers.google.com/p/$PROJECT_ID/r/$REPO', 'master']
```
