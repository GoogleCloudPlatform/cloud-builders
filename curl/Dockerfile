FROM launcher.gcr.io/google/ubuntu16_04

RUN apt-get update && \
  apt-get -y install curl ca-certificates

COPY deprecation.sh /usr/bin
ENTRYPOINT ["/usr/bin/deprecation.sh"]
