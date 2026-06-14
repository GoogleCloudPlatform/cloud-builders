FROM gcr.io/gcp-runtimes/ubuntu_20_0_4
LABEL org.opencontainers.image.source="https://github.com/GoogleCloudPlatform/cloud-builders" \
      org.opencontainers.image.vendor="Google LLC"

RUN apt-get update -qqy && apt-get dist-upgrade -yq && \
    apt-get -y install wget ca-certificates

COPY notice.sh /usr/bin
ENTRYPOINT ["/usr/bin/notice.sh"]
