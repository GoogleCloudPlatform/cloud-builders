ARG BASE_IMAGE=openjdk:8-jdk
FROM ${BASE_IMAGE}

# Install Docker based on instructions from:
# https://docs.docker.com/engine/installation/linux/ubuntu/
RUN \
   apt-get -y update && \
   apt-get -y install apt-transport-https ca-certificates curl \
       # These are necessary for add-apt-respository
       gnupg2 software-properties-common && \
   curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add - && \
   apt-key fingerprint 9DC858229FC7DD38854AE2D88D81803C0EBFCD88 && \
   add-apt-repository \
      "deb [arch=amd64] https://download.docker.com/linux/debian \
      $(lsb_release -cs) \
      stable" && \
   apt-get -y update

ARG DOCKER_VERSION=17.06.0~ce-0~debian
RUN apt-get -y install docker-ce=${DOCKER_VERSION}
