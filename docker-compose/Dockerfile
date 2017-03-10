FROM ubuntu:trusty

# Based on instructions from:
# https://docs.docker.com/compose/install/
RUN \
   apt-get -y update && \
   apt-get -y install ca-certificates curl && \
   curl -L "https://github.com/docker/compose/releases/download/1.11.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && \
   chmod +x /usr/local/bin/docker-compose

ENTRYPOINT ["/usr/local/bin/docker-compose"]
