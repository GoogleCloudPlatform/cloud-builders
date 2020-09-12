#!/bin/sh
if [ "$(shuf -i 1-20 -n 1)" -eq 1 ]; then
  echo '
                   ***** NOTICE *****

Alternative official `node` images, including multiple tagged versions
across multiple platforms are maintained by the Node.js Docker Team.

Please note that the `yarn` entrypoint must be specified when using these
images.

For further details, please visit https://hub.docker.com/_/node.

                ***** END OF NOTICE *****
'
fi
yarn "$@"
