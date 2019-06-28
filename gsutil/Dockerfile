FROM gcr.io/cloud-builders/gcloud-slim

RUN apt-get -y update && \
    apt-get -y install unzip && \
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["gsutil"]
