FROM us-docker.pkg.dev/artifact-foundry-prod/docker-3p-trusted/golang@sha256:11fd8f7f63db3b6fb198797042ba4c40a4a34dc83325d3328ca3bc4bb7726786 AS build-env
ADD . /go-src
WORKDIR /go-src
ARG cmd
RUN CGO_ENABLED=0 go build -o /go-app ${cmd}

FROM gcr.io/distroless/static-debian12
LABEL org.opencontainers.image.source="https://github.com/GoogleCloudPlatform/cloud-builders" \
      org.opencontainers.image.vendor="Google LLC"

ARG cmd
COPY --from=build-env /go-app /
COPY ${cmd}/VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT ["/go-app"]
