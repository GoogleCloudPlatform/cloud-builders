FROM gcr.io/cloud-builders/gcloud
COPY helm.bash /builder/helm.bash
RUN chmod +x /builder/helm.bash && mkdir -p /builder/helm && apt-get update && apt-get install -y curl \
    && curl -SL https://storage.googleapis.com/kubernetes-helm/helm-v2.5.1-linux-amd64.tar.gz -o helm.tar.gz \
    && tar zxvf helm.tar.gz --strip-components=1 -C /builder/helm linux-amd64/helm \
    && rm helm.tar.gz \
    && apt-get remove --purge -y curl && apt-get --purge -y autoremove && apt-get clean && rm -rf /var/lib/apt/lists/*
ENV PATH=/builder/helm/:$PATH
ENTRYPOINT ["/builder/helm.bash"]
