#!/bin/sh
if [[ $(( $RANDOM % 20 )) -eq 1 ]]; then
  echo '
                   ***** NOTICE *****

An official `microsoft/dotnet:sdk` image to run the `dotnet` tool exists at
https://hub.docker.com/r/microsoft/dotnet.

                ***** END OF NOTICE *****
'
fi
dotnet "$@"
