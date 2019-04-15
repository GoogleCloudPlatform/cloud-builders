FROM gcr.io/cloud-builders/gcloud-slim

RUN apt-get -y update && \
    apt-get -y install wget make python-setuptools && \
    # Setup Cloud Foundation Toolkit from GitHub (master)
    wget -qO- https://github.com/GoogleCloudPlatform/deploymentmanager-samples/tarball/master | tar zxv && \
    cd GoogleCloudPlatform-deploymentmanager-samples-*/community/cloud-foundation/ && \
    make clean && \
    make prerequisites && \
    make build && \
    make install && \
    # Clean up
    apt-get -y remove make python-setuptools && \
    rm -rf /var/lib/apt/lists/* && \
    cd / && \
    rm -rf /GoogleCloudPlatform-deploymentmanager-*/

ENTRYPOINT ["/usr/local/bin/cft"]
CMD ["--version"]
