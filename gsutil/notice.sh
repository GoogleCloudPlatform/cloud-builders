#!/bin/sh
if [[ $(( $RANDOM % 20 )) -eq 1 ]]; then
  echo '
                   ***** NOTICE *****

Official `cloud-sdk` images, including multiple tagged versions across multiple
platforms, can be found at
https://github.com/GoogleCloudPlatform/cloud-sdk-docker.

Suggested alternative images include:

    gcr.io/google.com/cloudsdktool/cloud-sdk
    gcr.io/google.com/cloudsdktool/cloud-sdk:alpine
    gcr.io/google.com/cloudsdktool/cloud-sdk:debian_component_based
    gcr.io/google.com/cloudsdktool/cloud-sdk:slim

Please note that the `gsutil` entrypoint must be specified when using these
images.

                ***** END OF NOTICE *****
'
fi
/builder/google-cloud-sdk/bin/gsutil "$@"
