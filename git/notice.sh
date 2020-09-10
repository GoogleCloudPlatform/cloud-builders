#!/bin/sh
if [ "$(shuf -i 1-20 -n 1)" -eq 1 ]; then
  echo '
                   ***** NOTICE *****

Official `cloud-sdk` images, including multiple tagged versions across multiple
platforms, can be found at
https://github.com/GoogleCloudPlatform/cloud-sdk-docker and may be more suitable
for some use cases when interacting with Cloud Source Repositories.

For additional information, please visit
https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/git

                ***** END OF NOTICE *****
'
fi
git "$@"
