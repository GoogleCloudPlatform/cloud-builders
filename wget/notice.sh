#!/bin/sh
if [ "$(shuf -i 1-20 -n 1)" -eq 1 ]; then
  echo '
                   ***** NOTICE *****

Please visit
https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/wget
where the README.md illustrates usage of numerous images that may provide better
alternatives to the use of wget.

                ***** END OF NOTICE *****
'
fi
wget "$@"
