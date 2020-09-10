#!/bin/sh
if [ "$(shuf -i 1-20 -n 1)" -eq 1 ]; then
  echo '
                   ***** NOTICE *****

Alternative official `gradle` images, including multiple tagged versions across
multiple platforms, can be found at https://hub.docker.com/_/gradle.

                ***** END OF NOTICE *****
'
fi
/usr/bin/gradle "$@"
