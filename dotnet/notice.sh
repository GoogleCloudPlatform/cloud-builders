#!/bin/sh
echo '
                   ***** NOTICE *****

An official `microsoft/dotnet:sdk` image to run the `dotnet` tool exists at
https://hub.docker.com/r/microsoft/dotnet.

                ***** END OF NOTICE *****
'
dotnet "$@"
