#!/bin/sh
echo '
                   ***** NOTICE *****

Alternative official `gradle` images, including multiple taged version across
multiple platforms, can be found at https://hub.docker.com/_/gradle.

                ***** END OF NOTICE *****
'

/usr/bin/gradle $@
