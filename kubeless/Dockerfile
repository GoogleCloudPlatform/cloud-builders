FROM gcr.io/cloud-builders/kubectl

RUN apt-get -y update && \
    apt-get -y install wget ca-certificates unzip

RUN mkdir -p /builder/bin && \
    wget -q https://github.com/kubeless/kubeless/releases/download/v1.0.0-alpha.7/kubeless_linux-amd64.zip && \
    unzip kubeless_linux-amd64.zip -d /builder && \
    mv /builder/bundles/kubeless_linux-amd64/kubeless /builder/bin/kubeless && \

    # Clean up
    apt-get -y remove wget && \
    rm -rf /var/lib/apt/lists/*

COPY kubeless.bash /builder/kubeless.bash

ENV PATH=/builder/bin:$PATH

ENTRYPOINT ["/builder/kubeless.bash"]

