#!/bin/sh
if [[ $(( $RANDOM % 20 )) -eq 1 ]]; then
  echo '
                   ***** NOTICE *****

Alternative official `gradle` images, including multiple tagged versions across
multiple platforms, can be found at https://hub.docker.com/_/gradle.

                ***** END OF NOTICE *****
'
fi
/usr/bin/gradle "$@"
