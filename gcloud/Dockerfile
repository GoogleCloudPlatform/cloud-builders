FROM ubuntu:trusty

RUN apt-get -y update && \
    apt-get -y install python2.7 wget unzip ca-certificates \
       # These are necessary for add-apt-respository
       software-properties-common python-software-properties && \
    # Install Git >2.0.1
    add-apt-repository ppa:git-core/ppa && \
    apt-get -y update && \
    apt-get -y install git && \
    # Setup Google Cloud SDK (latest)
    wget -nv https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && \
    mkdir -p /builder && \
    unzip -qq google-cloud-sdk.zip -d /builder && \
    rm google-cloud-sdk.zip && \
    CLOUDSDK_PYTHON="python2.7" /builder/google-cloud-sdk/install.sh --usage-reporting=false \
        --bash-completion=false \
        --disable-installation-options && \
    # Install alpha and beta components
    /builder/google-cloud-sdk/bin/gcloud -q components install alpha beta && \
    /builder/google-cloud-sdk/bin/gcloud -q components update && \

    # TODO(jasonhall): These lines pin gcloud to a particular version.
    # /builder/google-cloud-sdk/bin/gcloud components update --version 137.0.0 && \
    # /builder/google-cloud-sdk/bin/gcloud config set component_manager/disable_update_check 1 && \
    # /builder/google-cloud-sdk/bin/gcloud -q components update && \

    # Clean up
    rm -rf ~/.config/gcloud

ENV PATH=$PATH:/builder/google-cloud-sdk/bin/

ENTRYPOINT ["/builder/google-cloud-sdk/bin/gcloud"]
