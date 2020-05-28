#!/bin/sh

echo '
                   ***** NOTICE *****

Official `cloud-sdk` images, including multiple tagged versions across multiple
platforms, can be found at
https://github.com/GoogleCloudPlatform/cloud-sdk-docker.

Suggested alternative images include:

    gcr.io/google.com/cloudsdktool/cloud-sdk
    gcr.io/google.com/cloudsdktool/cloud-sdk:slim
    gcr.io/google.com/cloudsdktool/cloud-sdk:alpine

Please note that the `gsutil` entrypoint must be specified when using these
images.

                ***** END OF NOTICE *****
'

/builder/google-cloud-sdk/bin/gsutil "$@"
