#!/bin/sh

echo           \*\*\*\*\* DEPRECATION NOTICE \*\*\*\*\*
echo
echo This image is deprecated and will no longer be updated.
echo This recent version of the image will continue to exist.
echo
echo In place of this image, please use one of the following
echo images from https://hub.docker.com/r/google/cloud-sdk/:
echo
echo     google/cloud-sdk
echo     google/cloud-sdk:slim
echo     google/cloud-sdk:alpine
echo
echo To run \`gsutil\` with any of these images, you\'ll need to
echo specify the \`gsutil\` command as an argument or entrypoint.
echo Please note that these images support pinned versions
echo as well.
echo
echo           \*\*\*\*\* DEPRECATION NOTICE \*\*\*\*\*

/builder/google-cloud-sdk/bin/gsutil $@
