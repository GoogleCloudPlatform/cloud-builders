#!/bin/sh
echo '
                   ***** NOTICE *****

Alternative official `node` images, including multiple tagged versions
across multiple platforms are maintained by the Node.js Docker Team.

Please note that when using these images, you will need to specify the `yarn`
entrypoint.

For further details, please visit https://hub.docker.com/_/node.

                ***** END OF NOTICE *****
'

yarn $@
