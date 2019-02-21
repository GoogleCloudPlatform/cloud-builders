FROM google/cloud-sdk:alpine

# Install kubectl component
RUN /google-cloud-sdk/bin/gcloud -q components install kubectl

COPY kubectl.sh /builder/kubectl.sh

ENTRYPOINT ["/builder/kubectl.sh"]
