FROM golang:1.21 AS build-env
ADD . /go-src
WORKDIR /go-src
ARG cmd
RUN CGO_ENABLED=0 go build -o /go-app ${cmd}

FROM gcr.io/distroless/static-debian12
ARG cmd
COPY --from=build-env /go-app /
COPY ${cmd}/VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT ["/go-app"]
