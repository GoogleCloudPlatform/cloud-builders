FROM ubuntu:14.04

RUN apt-get update -qqy
RUN apt-get install -qqy curl
RUN curl https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 > /usr/bin/docker && chmod +x /usr/bin/docker

ENTRYPOINT ["/usr/bin/docker"]
