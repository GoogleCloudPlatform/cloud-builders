FROM ubuntu:trusty

# Based on instructions from:
# https://docs.docker.com/engine/installation/linux/ubuntu/
RUN \
   apt-get -y update && \
   apt-get -y install apt-transport-https ca-certificates curl \
       # These are necessary for add-apt-respository
       software-properties-common python-software-properties && \
   curl -fsSL https://yum.dockerproject.org/gpg | sudo apt-key add - && \
   apt-key fingerprint 58118E89F3A912897C070ADBF76221572C52609D && \
   add-apt-repository \
       "deb https://apt.dockerproject.org/repo/ \
       ubuntu-$(lsb_release -cs) \
       main" && \
   apt-get -y update

ARG DOCKER_VERSION=17.05.0~ce-0~ubuntu-trusty
RUN apt-get -y install docker-engine=${DOCKER_VERSION}

ENTRYPOINT ["/usr/bin/docker"]
