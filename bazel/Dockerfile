FROM gcr.io/gcp-runtimes/ubuntu_20_0_4

ADD bazel.sh /builder/bazel.sh
ARG BAZEL_VERSION=latest
ARG DOCKER_VERSION=5:24.0.7-1~ubuntu.20.04~focal

RUN \
    # This makes add-apt-repository available.
    apt-get update && \
    apt-get -y install \
        python \
        python3 \
        python-pkg-resources \
        python3-pkg-resources \
        software-properties-common \
        unzip && \

    # Install Git >2.0.1
    add-apt-repository ppa:git-core/ppa && \
    apt-get -y update && \
    apt-get -y install git && \

    # Install bazel (https://docs.bazel.build/versions/master/install-ubuntu.html)
    apt-get -y install openjdk-8-jdk && \
    echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list && \
    curl https://bazel.build/bazel-release.pub.gpg | apt-key add - && \
    apt-get update && \

    apt-get -y install bazel-${BAZEL_VERSION} && \
    apt-get -y upgrade bazel-${BAZEL_VERSION} && \

    ln -s /usr/bin/bazel-${BAZEL_VERSION} /usr/bin/bazel && \

    # Install Docker (https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/#uninstall-old-versions)
    apt-get -y install \
        apt-transport-https \
        curl \
        ca-certificates && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository \
      "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) \
      stable edge" && \
    apt-get -y update && \
    apt-get dist-upgrade -y && \
    apt-get install -y docker-ce=${DOCKER_VERSION} docker-ce-cli=${DOCKER_VERSION} unzip && \
    rm -rf \
       /var/cache/debconf/* \
       /var/lib/apt/lists/* \
       /var/log/* \
       /tmp/* \
       /var/tmp/*  && \

    mv /usr/bin/bazel /builder/bazel           && \
    mv /builder/bazel.sh /usr/bin/bazel        && \

    # Unpack bazel for future use.
    bazel version

# Store the Bazel outputs under /workspace so that the symlinks under bazel-bin (et al) are accessible
# to downstream build steps.
RUN mkdir -p /workspace
RUN echo 'startup --output_base=/workspace/.bazel' > ~/.bazelrc

ENTRYPOINT ["bazel"]
