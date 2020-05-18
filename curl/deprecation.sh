#!/bin/sh
echo
echo           \*\*\*\*\* DEPRECATION NOTICE \*\*\*\*\*
echo
echo This image is deprecated and will no longer be updated.
echo This recent version of the image will continue to exist.
echo
echo For best support of \`curl\` please use one of the official
echo \`curl\` images maintained by the curlimages community
echo on Dockerhub.
echo For details, visit https://hub.docker.com/r/curlimages/curl. Note that
echo this image executes as special user \`curl_user\` and thus may not be
echo suitable for all purposes.
echo
echo Alternatively, image \`launcher.gcr.io/google/ubuntu1604\` is maintained
echo by Google, has \`curl\` installed, and executes as \`root\`.
echo For details, visit
echo https://console.cloud.google.com/launcher/details/google/ubuntu1604.
echo
echo           \*\*\*\*\* DEPRECATION NOTICE \*\*\*\*\*
echo
curl $@
