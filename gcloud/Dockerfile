FROM python:2.7-alpine

RUN apk update && \
    apk add wget unzip ca-certificates git && \
    # Setup Google Cloud SDK (latest)
    wget -nv https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && \
    mkdir -p /builder && \
    unzip -qq google-cloud-sdk.zip -d /builder && \
    rm google-cloud-sdk.zip && \
    /builder/google-cloud-sdk/install.sh --usage-reporting=false \
        --bash-completion=false \
        --disable-installation-options && \
    # Install alpha and beta components
    /builder/google-cloud-sdk/bin/gcloud -q components install alpha beta && \
    /builder/google-cloud-sdk/bin/gcloud -q components update && \

    # TODO(jasonhall): These lines pin gcloud to a particular version.
    # google-cloud-sdk/bin/gcloud components update --version 110.0.0 && \
    # google-cloud-sdk/bin/gcloud config set component_manager/disable_update_check 1 && \
    # google-cloud-sdk/bin/gcloud -q components update && \

    # Clean up
    apk del wget unzip && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["/builder/google-cloud-sdk/bin/gcloud"]
