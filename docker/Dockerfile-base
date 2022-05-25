FROM gcr.io/gcp-runtimes/ubuntu_20_0_4

RUN apt-get -y update && \
    apt-get -y install \
        apt-transport-https \
        ca-certificates \
        curl \
        make \
        software-properties-common && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | \
        gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) \
        signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
        https://download.docker.com/linux/ubuntu \
        $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list && \
    apt-get -y update && \
    apt-get -y dist-upgrade && \
    apt-get autoremove && \
    apt-get clean
