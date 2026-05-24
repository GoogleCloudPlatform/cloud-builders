FROM gcr.io/gcp-runtimes/ubuntu_20_0_4
LABEL org.opencontainers.image.source="https://github.com/GoogleCloudPlatform/cloud-builders" \
      org.opencontainers.image.vendor="Google LLC"

RUN apt-get update && apt-get dist-upgrade -y && \
    rm -rf \
       /var/cache/debconf/* \
       /var/lib/apt/lists/* \
       /var/log/* \
       /tmp/* \
       /var/tmp/*
COPY notice.sh /usr/bin
ENTRYPOINT ["/usr/bin/notice.sh"]
