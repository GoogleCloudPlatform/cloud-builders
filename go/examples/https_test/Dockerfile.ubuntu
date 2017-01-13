FROM ubuntu

RUN apt-get update -qqy && apt-get install -qqy ca-certificates

COPY gopath/bin/https_test /https_test

ENTRYPOINT ["/https_test"]
