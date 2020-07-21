#!/bin/sh
if [[ $(( $RANDOM % 20 )) -eq 1 ]]; then
  echo '
                   ***** NOTICE *****

Alternative official `maven` images, including multiple tagged versions
across multiple jdk versions and multiple platforms, can be found at
https://hub.docker.com/_/maven.

                ***** END OF NOTICE *****
'
fi
mvn "$@"
