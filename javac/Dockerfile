ARG BASE_IMAGE=launcher.gcr.io/google/openjdk8
FROM ${BASE_IMAGE}

ARG DOCKER_VERSION=5:19.03.8~3-0~debian-stretch
ARG PYTHON_VERSION=3.11.5

# https://lists.debian.org/debian-devel-announce/2023/03/msg00006.html
# Base image uses debian9/stretch, and as such needs to have
# sources.list point to archive.debian.org.
RUN \
   sed -i s/httpredir.debian.org/archive.debian.org/g /etc/apt/sources.list && \
   sed -i s/deb.debian.org/archive.debian.org/g /etc/apt/sources.list && \
   sed -i 's|security.debian.org|archive.debian.org/debian-security/|g' /etc/apt/sources.list && \
   sed -i '/stretch-updates/d' /etc/apt/sources.list

# Install Docker based on instructions from:
# https://docs.docker.com/engine/installation/linux/docker-ce/debian
RUN \
   apt-get update -qqy && apt-get dist-upgrade -yq && \
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

# Need to install python from source because Debian 9 (stretch) doesn't provide supported versions of python (3.8+)
RUN \
   apt-get -y install python3-dev zlib1g-dev && \
   wget https://www.python.org/ftp/python/${PYTHON_VERSION}/Python-${PYTHON_VERSION}.tgz && \
   tar xzf Python-${PYTHON_VERSION}.tgz && \
   cd Python-${PYTHON_VERSION} && \
   ./configure --enable-optimizations && \
   make install

ENTRYPOINT ["javac"]
