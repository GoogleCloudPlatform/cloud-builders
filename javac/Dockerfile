ARG BASE_IMAGE=launcher.gcr.io/google/openjdk8
FROM ${BASE_IMAGE}

ARG DOCKER_VERSION=5:19.03.8~3-0~debian-stretch

# Install Docker based on instructions from:
# https://docs.docker.com/engine/installation/linux/docker-ce/debian
RUN \
   apt-get -y update && \
   apt-get --fix-broken -y install && \
   apt-get -y install apt-transport-https ca-certificates curl gnupg2 software-properties-common && \
   curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add - && \
   apt-key fingerprint 9DC858229FC7DD38854AE2D88D81803C0EBFCD88 && \
   add-apt-repository \
      "deb [arch=amd64] https://download.docker.com/linux/debian \
      $(lsb_release -cs) \
      stable" && \
   apt-get -y update && \
   apt-get -y install docker-ce=${DOCKER_VERSION} docker-ce-cli=${DOCKER_VERSION} && \

   # Clean up build packages
   apt-get remove -y --purge curl gnupg2 software-properties-common && \
   apt-get clean

ENTRYPOINT ["javac"]
