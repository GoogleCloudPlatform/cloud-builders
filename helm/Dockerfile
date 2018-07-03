FROM gcr.io/cloud-builders/kubectl

RUN curl https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz | tar -zx linux-amd64/helm && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    chmod +x /usr/local/bin/helm

COPY helm.bash /builder/helm.bash
RUN chmod +x /builder/helm.bash

ENTRYPOINT ["/builder/helm.bash"]
