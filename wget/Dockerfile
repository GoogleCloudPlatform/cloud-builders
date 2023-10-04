FROM gcr.io/gcp-runtimes/ubuntu_20_0_4

RUN apt-get update -qqy && apt-get dist-upgrade -yq && \
    apt-get -y install wget ca-certificates

COPY notice.sh /usr/bin
ENTRYPOINT ["/usr/bin/notice.sh"]
