FROM gcr.io/cloud-builders/gcloud-slim

RUN apt-get -y update && \
    apt-get -y install unzip zip && \
    rm -rf /var/lib/apt/lists/*

COPY notice.sh /builder

ENTRYPOINT ["/builder/notice.sh"]
