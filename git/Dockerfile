FROM gcr.io/cloud-builders/gcloud

ENV PATH=$PATH:/builder/google-cloud-sdk/bin/

RUN git config --system credential.helper gcloud.sh

ENTRYPOINT ["git"]
