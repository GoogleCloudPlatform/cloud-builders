FROM gcr.io/cloud-builders/gcloud-slim
LABEL org.opencontainers.image.source="https://github.com/GoogleCloudPlatform/cloud-builders" \
      org.opencontainers.image.vendor="Google LLC"

RUN apt-get -y update && \
    apt-get -y install unzip zip && \
    rm -rf /var/lib/apt/lists/*

COPY notice.sh /builder

ENTRYPOINT ["/builder/notice.sh"]
