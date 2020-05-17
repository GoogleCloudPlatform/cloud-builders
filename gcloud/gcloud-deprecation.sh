#!/bin/sh

echo           \*\*\*\*\* DEPRECATION NOTICE \*\*\*\*\*
echo
echo This image is deprecated and will no longer be updated.
echo This recent version of the image will continue to exist.
echo
echo In place of this image, please use one of the following
echo images built from
echo https://github.com/GoogleCloudPlatform/cloud-sdk-docker:
echo
echo     gcr.io/google.com/cloudsdktool/cloud-sdk
echo     gcr.io/google.com/cloudsdktool/cloud-sdk:slim
echo     gcr.io/google.com/cloudsdktool/cloud-sdk:alpine
echo
echo Please note that these images support pinned versions
echo as well.
echo
echo           \*\*\*\*\* DEPRECATION NOTICE \*\*\*\*\*

/builder/google-cloud-sdk/bin/gcloud $@
