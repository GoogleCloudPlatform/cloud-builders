FROM golang:buster AS build-env
ADD . /go-src
WORKDIR /go-src
ARG cmd
RUN go build -o /go-app ${cmd}

FROM gcr.io/distroless/base
ARG cmd
COPY --from=build-env /go-app /
COPY ${cmd}/VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT ["/go-app"]
