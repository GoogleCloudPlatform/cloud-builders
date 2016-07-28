FROM ubuntu:14.04

RUN apt-get update && \
    apt-get install -y python curl && \

    curl https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 > /usr/bin/docker && chmod +x /usr/bin/docker && \

    # This makes add-apt-repository available.
    apt-get -y install software-properties-common && \

    add-apt-repository ppa:webupd8team/java && \
    apt-get update && \

    # Accept the installer license && \
    echo "oracle-java7-installer shared/accepted-oracle-license-v1-1 boolean true" && \
    echo debconf shared/accepted-oracle-license-v1-1 select true | debconf-set-selections && \

    apt-get -y install oracle-java8-installer && \
    echo "deb http://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list && \
    apt-get -y install curl && \
    curl https://storage.googleapis.com/bazel-apt/doc/apt-key.pub.gpg | apt-key add - && \
    apt-get update && \

    apt-get -y install bazel && \

    # Unpack bazel for future use.
    bazel version

ENTRYPOINT ["bazel"]
