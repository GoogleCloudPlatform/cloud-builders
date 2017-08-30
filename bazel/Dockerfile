FROM ubuntu:14.04

RUN \
    # This makes add-apt-repository available.
    apt-get update && \
    apt-get -y install python software-properties-common unzip && \

    # Install Git >2.0.1
    add-apt-repository ppa:git-core/ppa && \
    apt-get -y update && \
    apt-get -y install git && \

    # Install Docker (https://docs.docker.com/engine/installation/linux/ubuntu/)
    apt-get -y install curl \
        apt-transport-https \
        ca-certificates && \
    curl -fsSL https://yum.dockerproject.org/gpg | sudo apt-key add - && \
    apt-key fingerprint 58118E89F3A912897C070ADBF76221572C52609D && \
    add-apt-repository \
       "deb https://apt.dockerproject.org/repo/ \
       ubuntu-$(lsb_release -cs) \
       main" && \
    apt-get update && \
    apt-get -y install docker-engine=17.05.0~ce-0~ubuntu-trusty && \
    add-apt-repository ppa:webupd8team/java && \
    apt-get update && \

    # Accept the installer license \
    echo "oracle-java7-installer shared/accepted-oracle-license-v1-1 boolean true" && \
    echo debconf shared/accepted-oracle-license-v1-1 select true | debconf-set-selections && \

    apt-get -y install oracle-java8-installer && \
    echo "deb http://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list && \
    curl https://storage.googleapis.com/bazel-apt/doc/apt-key.pub.gpg | apt-key add - && \
    apt-get update && \

    apt-get -y install bazel && \

    # Unpack bazel for future use.
    bazel version

# Store the Bazel outputs under /workspace so that the symlinks under bazel-bin (et al) are accessible
# to downstream build steps.
RUN echo 'startup --output_base=/workspace/.bazel' > ~/.bazelrc

ENTRYPOINT ["bazel"]
