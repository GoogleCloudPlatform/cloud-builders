FROM gcr.io/cloud-builders/gcloud-slim:ubuntu1804

# Install kubectl component
RUN /builder/google-cloud-sdk/bin/gcloud -q components install kubectl

COPY kubectl.bash /builder/kubectl.bash

ENTRYPOINT ["/builder/kubectl.bash"]
