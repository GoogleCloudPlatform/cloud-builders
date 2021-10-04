FROM gcr.io/cloud-builders/gcloud:ubuntu1804

RUN apt-get update && \
  apt-get -y install wget ca-certificates

COPY notice.sh /usr/bin
ENTRYPOINT ["/usr/bin/notice.sh"]
