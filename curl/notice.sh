#!/bin/sh
if [ "$(shuf -i 1-20 -n 1)" -eq 1 ]; then
  echo '
                   ***** NOTICE *****

Supported `curl` images, including multiple tagged versions,
are available at
https://console.cloud.google.com/launcher/details/google/ubuntu1604
and https://hub.docker.com/r/curlimages/curl.

                ***** END OF NOTICE *****
'
fi
curl "$@"
