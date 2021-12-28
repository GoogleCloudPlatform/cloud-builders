FROM gcloud-slim

RUN apt-get -y update && \
    # JRE is required for cloud-datastore-emulator
    apt-get -y install default-jre && \

    # Install all available components
    /builder/google-cloud-sdk/bin/gcloud -q components install \
        alpha beta \
        app-engine-go \
        app-engine-java \
        app-engine-python \
        app-engine-python-extras \
        bigtable \
        cbt \
        cloud-datastore-emulator \
        cloud-firestore-emulator \
        cloud-build-local \
        datalab \
        docker-credential-gcr \
        kpt \
        kubectl \
        kustomize \
        local-extract \
        pubsub-emulator \
        skaffold \
        && \

    /builder/google-cloud-sdk/bin/gcloud -q components update && \
    /builder/google-cloud-sdk/bin/gcloud components list && \

    # Clean up
    rm -rf /var/lib/apt/lists/*
