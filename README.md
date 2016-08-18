# Google Cloud Container Builder official builder images

This repository contains builder images for use with the [Google Cloud Container
Builder API](https://cloud.google.com/container-builder/docs/).

These images are available at `gcr.io/cloud-builders/...` and include such
images as:

*   `bazel`: runs the [bazel](https://bazel.io) tool
*   `docker`: runs docker commands directly
*   `dockerizer`: runs `docker build`
*   `git`: runs the [git](https://git-scm.com/) tool
*   `go`: runs the [go](https://golang.org/cmd/go) tool
*   `golang-project`: recognizes and builds conventional [Go](https://golang.org) projects into container images
*   `retagger`: re-tags an image and pushes to the repository
