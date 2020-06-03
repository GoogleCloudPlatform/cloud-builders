FROM launcher.gcr.io/google/ubuntu16_04

RUN apt-get update && \
  apt-get -y install wget ca-certificates

COPY notice.sh /usr/bin
ENTRYPOINT ["/usr/bin/notice.sh"]
