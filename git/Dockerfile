FROM gcr.io/cloud-builders/gcloud-slim
LABEL org.opencontainers.image.source="https://github.com/GoogleCloudPlatform/cloud-builders" \
      org.opencontainers.image.vendor="Google LLC"

COPY notice.sh /usr/bin
ENTRYPOINT ["/usr/bin/notice.sh"]
