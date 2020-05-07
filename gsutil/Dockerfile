FROM gcr.io/cloud-builders/gcloud-slim

RUN apt-get -y update && \
    apt-get -y install unzip zip && \
    rm -rf /var/lib/apt/lists/*

COPY ./gsutil-deprecation.sh /builder

ENTRYPOINT ["/builder/gsutil-deprecation.sh"]
