FROM gcloud-slim

RUN \
    # Install all available components
    /builder/google-cloud-sdk/bin/gcloud -q components install \
        alpha beta \
        app-engine-go \
        app-engine-java \
        app-engine-php \
        app-engine-python \
        app-engine-python-extras \
        bigtable \
        cbt \
        cloud-datastore-emulator \
        cloud-build-local \
        container-builder-local \
        datalab \
        docker-credential-gcr \
        emulator-reverse-proxy \
        kubectl \
        pubsub-emulator \
        && \

    /builder/google-cloud-sdk/bin/gcloud -q components update && \
    /builder/google-cloud-sdk/bin/gcloud components list
