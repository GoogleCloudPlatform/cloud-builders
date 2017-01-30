FROM alpine

# Docker installation instructions from:
# http://wiki.alpinelinux.org/wiki/Docker#Installation
RUN echo "http://nl.alpinelinux.org/alpine/v3.5/main" >> /etc/apk/repositories && \
    apk update && \
    apk add docker=1.12.6-r0

ENTRYPOINT ["/usr/bin/docker"]
