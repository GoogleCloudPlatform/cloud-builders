FROM alpine

# TODO: Upgrade docker in step with that used in the Container Builder service.
RUN apk update && \
    apk add curl && \
    curl https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 > /usr/bin/docker && chmod +x /usr/bin/docker && \
    apk del curl && \
    rm -rf /var/cache/apk/*

COPY dockerizer.ash /dockerizer.ash

ENTRYPOINT ["/dockerizer.ash"]
