FROM launcher.gcr.io/google/ubuntu16_04

RUN apt-get -y update && \
    apt-get -y install gcc python2.7 python-dev python-setuptools wget ca-certificates \
       # These are necessary for add-apt-respository
       software-properties-common python-software-properties \
       # JRE is required for cloud-datastore-emulator
       default-jre && \

    # Install Git >2.0.1
    add-apt-repository ppa:git-core/ppa && \
    apt-get -y update && \
    apt-get -y install git && \

    # Setup Google Cloud SDK (latest)
    mkdir -p /builder && \
    wget -qO- https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz | tar zxv -C /builder && \
    CLOUDSDK_PYTHON="python2.7" /builder/google-cloud-sdk/install.sh --usage-reporting=false \
        --bash-completion=false \
        --disable-installation-options && \

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
        container-builder-local \
        datalab \
        docker-credential-gcr \
        emulator-reverse-proxy \
        kubectl \
        pubsub-emulator \
        && \

    /builder/google-cloud-sdk/bin/gcloud -q components update && \
    /builder/google-cloud-sdk/bin/gcloud components list && \

    # install crcmod: https://cloud.google.com/storage/docs/gsutil/addlhelp/CRC32CandInstallingcrcmod
    easy_install -U pip && \
    pip install -U crcmod && \

    # TODO(jasonhall): These lines pin gcloud to a particular version.
    # /builder/google-cloud-sdk/bin/gcloud components update --version 137.0.0 && \
    # /builder/google-cloud-sdk/bin/gcloud config set component_manager/disable_update_check 1 && \
    # /builder/google-cloud-sdk/bin/gcloud -q components update && \

    # Clean up
    apt-get -y remove gcc python-dev python-setuptools wget && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf ~/.config/gcloud

ENV PATH=/builder/google-cloud-sdk/bin/:$PATH

ENTRYPOINT ["gcloud"]
