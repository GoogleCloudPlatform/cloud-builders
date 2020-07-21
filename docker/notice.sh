#!/bin/sh
if [[ $(( $RANDOM % 20 )) -eq 1 ]]; then
  echo '
                   ***** NOTICE *****

Alternative official `docker` images, including multiple versions across
multiple platforms, are maintained by the Docker Team. For details, please
visit https://hub.docker.com/_/docker.

                ***** END OF NOTICE *****
'
fi
/usr/bin/docker "$@"
