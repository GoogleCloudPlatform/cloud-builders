# Google Cloud Container Builder official builder images

This repository contains builder images for use with the [Google Cloud Container
Builder API](https://cloud.google.com/container-builder/docs/).

These images are available at `gcr.io/cloud-builders/...` and include such
images as:

*   `docker`: runs docker commands directly
*   `dockerizer`: runs `docker build`
*   `retagger`: re-tags an image and pushes to the repository
