FROM gcr.io/cloud-builders/gcloud-slim
LABEL org.opencontainers.image.source="https://github.com/GoogleCloudPlatform/cloud-builders" \
      org.opencontainers.image.vendor="Google LLC"

# Install kubectl component
RUN /builder/google-cloud-sdk/bin/gcloud -q components install kubectl

COPY kubectl.bash /builder/kubectl.bash

ENTRYPOINT ["/builder/kubectl.bash"]
