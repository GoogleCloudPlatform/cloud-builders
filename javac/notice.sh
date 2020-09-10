#!/bin/sh
if [ "$(shuf -i 1-20 -n 1)" -eq 1 ]; then
  echo '
                   ***** NOTICE *****

Alternative official `openjdk` images, including multiple tagged versions
across multiple jdk versions and multiple platforms, can be found at
https://hub.docker.com/_/openjdk.

                ***** END OF NOTICE *****
'
fi
javac $@
