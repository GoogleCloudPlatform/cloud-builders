# Google Cloud Build official builder images

This repository contains source code for official builders used with the [Google
Cloud Build API](https://cloud.google.com/cloud-build/docs/).

Pre-built images are available at `gcr.io/cloud-builders/...` and include:

*   `bazel`: runs the [bazel](https://bazel.io) tool
*   `curl`: runs the [curl](https://curl.haxx.se) tool
*   `docker`: runs the [docker](https://docker.com) tool
*   `dotnet`: run the [dotnet](https://docs.microsoft.com/dotnet/core/tools/) tool
*   `gcloud`: runs the [gcloud](https://cloud.google.com/sdk/gcloud/) tool
*   `gcs-fetcher`: efficiently fetches objects from Google Cloud Storage
*   `git`: runs the [git](https://git-scm.com/) tool
*   `gke-deploy`: deploys an application to a Kubernetes cluster, following Google's recommended best practices
*   `go`: runs the [go](https://golang.org/cmd/go) tool
*   `gradle`: runs the [gradle](https://gradle.org/) tool
*   `gsutil`: runs the [gsutil](https://cloud.google.com/storage/docs/gsutil) tool
*   `javac`: runs the [javac](https://docs.oracle.com/javase/7/docs/technotes/tools/windows/javac.html) tool
*   `kubectl`: runs the [kubectl](https://kubernetes.io/docs/user-guide/kubectl-overview/) tool
*   `mvn`: runs the [maven](https://maven.apache.org/) tool
*   `npm`: runs the [npm](https://docs.npmjs.com/) tool
*   `twine`: runs the [twine](https://twine.readthedocs.io/) tool
*   `wget`: runs the [wget](https://www.gnu.org/software/wget/) tool
*   `yarn`: runs the [yarn](https://yarnpkg.com/) tool

Builders contributed by the public are available in the [Cloud Builders
Community
repo](https://github.com/GoogleCloudPlatform/cloud-builders-community).

Each builder includes a `cloudbuild.yaml` that will push your images to [Artifact
Registry](https://cloud.google.com/artifact-registry) in addition to [Google Container
Registry](https://cloud.google.com/container-registry). To build with this default `cloudbuild.yaml`,
you will need to first [create a Docker
repository](https://cloud.google.com/artifact-registry/docs/docker/store-docker-container-images#create)
to store the images. The provided `cloudbuild.yaml` assumes your project has set up a [multi-region](https://cloud.google.com/artifact-registry/docs/repositories/repo-locations#location-mr) Artifact
Registry Docker repositories called `ga` and that is setup for `us`, `europe`, and `asia` multi-regions (i.e. `us-docker.pkg.dev`, `europe-docker.pkg.dev`, `asia-docker.pkg.dev`).

To file issues and feature requests against these builder images,
[create an issue in this repo](https://github.com/GoogleCloudPlatform/cloud-builders/issues/new).
If you are experiencing an issue with the Cloud Build service or
have a feature request, e-mail google-cloud-dev@googlegroups.com
or see our [Getting support](https://cloud.google.com/cloud-build/docs/getting-support)
documentation.

---

# Alternative Official Images

Most of the tools in this repo are also available in official
community-supported publicly available repositories. Such
repos also generally support multiple versions and platforms,
available by tag.

The following official community-supported images are compatible with the
hosted Cloud Build service and function well as build steps; note that
some will require that you specify an `entrypoint` for the image. Additional
details regarding each alternative official image are available in the `README.md`
for the corresponding Cloud Builder.

*   [`docker`](https://hub.docker.com/_/docker/) supports tagged docker versions across multiple platforms
*   [`gcr.io/google.com/cloudsdktool/cloud-sdk`](https://github.com/GoogleCloudPlatform/cloud-sdk-docker) includes multiple entrypoints:
    *   `gcloud`: runs the [gcloud](https://cloud.google.com/sdk/gcloud/) tool
    *   `gsutil`: runs the [gsutil](https://cloud.google.com/storage/docs/gsutil) tool
    *   `kubectl`: runs the [kubectl](https://kubernetes.io/docs/user-guide/kubectl-overview/) tool
*   [`node`](https://hub.docker.com/_/node) includes these entrypoints:
    *   `npm`: runs the [npm](https://docs.npmjs.com/) tool
    *   `yarn`: runs the [yarn](https://yarnpkg.com/) tool
*   [`microsoft/dotnet:sdk`](https://hub.docker.com/_/microsoft-dotnet-core) includes
    *   `dotnet`: runs the [dotnet](https://docs.microsoft.com/dotnet/core/tools/) tool
*   Java builders include:
    *   [`openjdk`](https://hub.docker.com/_/openjdk) supports many production versions of Java across multiple platforms
    *   [`gradle`](https://hub.docker.com/_/gradle/) supports a matrix of Java and gradle versions across multiple platforms
    *   [`maven`](https://hub.docker.com/_/maven/) supports a matrix of Java and maven versions across multiple platforms
*   [`gcr.io/cloud-marketplace-containers/google/bazel`](http://gcr.io/cloud-marketplace-containers/google/bazel) is provided by the bazel team and runs the [`bazel`](https://bazel.build/) tool
*   `curl` is packaged in:
    *   [`launcher.gcr.io/google/ubuntu1604`](https://console.cloud.google.com/launcher/details/google/ubuntu1604)
    *   [`curlimages/curl`](https://hub.docker.com/r/curlimages/curl) is community-supported
*   [`golang`](https://hub.docker.com/_/golang) is provided by the Go team and runs the [`go`](https://golang.org/cmd/go/) tool

# Future Direction

You may have already noticed that most of the images in this repo now provide notices to the
above alternative images. For the hosted Cloud Build service, we are formulating plans
surrounding both improved support for existing `cloud-builder` images and documentation for
alternative community-supported images that may be more appropriate for some users. Both this
page and the related [open issues](https://github.com/GoogleCloudPlatform/cloud-builders/labels/augmentation)
will be updated with details soon.
