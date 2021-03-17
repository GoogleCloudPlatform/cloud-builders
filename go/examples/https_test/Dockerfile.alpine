FROM alpine

RUN apk add --update ca-certificates

COPY gopath/bin/https_test /https_test

ENTRYPOINT ["/https_test"]
