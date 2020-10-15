#!/bin/sh
if [ "$(shuf -i 1-20 -n 1)" -eq 1 ]; then
  echo '
                   ***** NOTICE *****

Supported `curl` versions can be found in the various images available at
https://console.cloud.google.com/launcher/details/google/ubuntu1604.

                ***** END OF NOTICE *****
'
fi
curl "$@"
