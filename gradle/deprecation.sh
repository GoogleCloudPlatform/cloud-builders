#!/bin/sh
echo
echo           \*\*\*\*\* DEPRECATION NOTICE \*\*\*\*\*
echo
echo This image is deprecated and will no longer be updated.
echo This recent version of the image will continue to exist.
echo
echo In place of this image, please use one of the official
echo \`gradle\` images.
echo For details, visit https://hub.docker.com/_/gradle.
echo
echo           \*\*\*\*\* DEPRECATION NOTICE \*\*\*\*\*
echo
/usr/bin/gradle $@
