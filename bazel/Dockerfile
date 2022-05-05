FROM gcr.io/gcp-runtimes/ubuntu_20_0_4

ADD bazel.sh /builder/bazel.sh
ARG DOCKER_VERSION=5:19.03.9~3-0~ubuntu-focal

RUN \
    # This makes add-apt-repository available.
    apt-get update && \
    apt-get -y install \
        apt-transport-https \
        curl \
        gnupg2 \
        python \
        python-pkg-resources \
        python3 \
        python3-pkg-resources \
        software-properties-common \
        unzip && \

    # Install Git >2.0.1
    add-apt-repository ppa:git-core/ppa && \
    apt-get -y update && \
    apt-get -y install git && \

    # Install latest bazel (https://docs.bazel.build/versions/master/install-ubuntu.html)
    curl -fsSL https://bazel.build/bazel-release.pub.gpg | apt-key add - && \
    echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list && \
    apt-get update && \
    apt-get -y install bazel && \
    apt-get -y upgrade bazel && \

    # Install Docker (https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/#uninstall-old-versions)
    apt-get -y install \
        linux-image-extra-virtual \
        apt-transport-https \
        curl \
        ca-certificates && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository \
      "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) \
      stable edge" && \
    apt-get -y update && \
    apt-get install -y docker-ce=${DOCKER_VERSION} docker-ce-cli=${DOCKER_VERSION} unzip && \

    # Install JDK (JDK 11 is default for ubuntu 20.0.4) for building java code
    apt-get install -y default-jdk && \

    # Replace the default bazel binary with script that adds Cloud Build customizations
    mv /usr/bin/bazel /builder/bazel           && \
    mv /usr/bin/bazel-real /builder/bazel-real && \
    mv /builder/bazel.sh /usr/bin/bazel        && \

    # Store the Bazel outputs under /workspace so that the symlinks under bazel-bin (et al) are accessible
    # to downstream build steps.
    mkdir -p /workspace && \
    echo 'startup --output_base=/workspace/.bazel' >> /etc/bazel.bazelrc

ENTRYPOINT ["bazel"]
