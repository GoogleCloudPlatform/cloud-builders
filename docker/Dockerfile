FROM alpine

RUN apk update && \
    apk add curl && \
    curl https://get.docker.com/builds/Linux/x86_64/docker-1.12.6 > /usr/bin/docker && chmod +x /usr/bin/docker && \
    apk del curl && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/bin/docker"]
