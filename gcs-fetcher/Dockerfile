FROM golang:stretch AS build-env
ADD . /go-src
WORKDIR /go-src
ARG cmd
RUN go build -o /go-app ${cmd}

FROM gcr.io/distroless/base
COPY --from=build-env /go-app /
ENTRYPOINT ["/go-app"] 
