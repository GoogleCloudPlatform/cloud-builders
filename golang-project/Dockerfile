FROM gcr.io/cloud-builders/go

# TODO: Upgrade docker in step with that used in the Container Builder service.
RUN apk update && \
    apk upgrade && \
    apk add curl && \
    curl https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 > /usr/bin/docker && chmod +x /usr/bin/docker && \
    apk del curl && \
    rm -rf /var/cache/apk/*

COPY gopath/bin/golang_project /builder/golang_project
COPY golang_project.sh /builder/golang_project.sh

ENTRYPOINT ["/builder/golang_project.sh"]
