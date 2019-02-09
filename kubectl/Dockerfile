FROM gcr.io/cloud-builders/gcloud-slim

# Install kubectl. gcloud's component manager is disabled, so we have to
# install it from apt-get.
RUN apt-get install kubectl

COPY kubectl.bash /builder/kubectl.bash

ENTRYPOINT ["/builder/kubectl.bash"]
