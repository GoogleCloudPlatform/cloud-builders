FROM launcher.gcr.io/google/ubuntu16_04

RUN apt-get -y update && \
    apt-get -y install \
        apt-transport-https \
        ca-certificates \
        curl \
        make \
        software-properties-common && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    apt-key fingerprint 0EBFCD88 && \
    add-apt-repository \
       "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
       xenial \
       stable" && \
    apt-get -y update && \
    apt-get -y install docker-ce docker-ce-cli

ENTRYPOINT ["/usr/bin/docker"]
