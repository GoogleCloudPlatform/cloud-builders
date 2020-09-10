#!/bin/sh
if [ "$(shuf -i 1-20 -n 1)" -eq 1 ]; then
  echo '
                   ***** NOTICE *****

An official `microsoft/dotnet:sdk` image to run the `dotnet` tool exists at
https://hub.docker.com/r/microsoft/dotnet.

                ***** END OF NOTICE *****
'
fi
dotnet "$@"
