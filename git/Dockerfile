FROM gcr.io/cloud-builders/gcloud

RUN git config --system credential.helper gcloud.sh

COPY notice.sh /usr/bin
ENTRYPOINT ["/usr/bin/notice.sh"]
