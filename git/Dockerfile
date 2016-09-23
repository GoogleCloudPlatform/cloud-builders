FROM gcr.io/cloud-builders/gcloud

ENV PATH=$PATH:/builder/google-cloud-sdk/bin/

RUN apk update && \
  apk add git && \
  rm -rf /var/cache/apk/* && \

  # Configure the gcloud credential helper to provide builder service account
  # credentials.
  git config --global credential.helper gcloud.sh

ENTRYPOINT ["git"]
