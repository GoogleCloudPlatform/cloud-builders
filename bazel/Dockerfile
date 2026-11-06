FROM gcr.io/gcp-runtimes/ubuntu_22_0_4

ADD bazel.sh /builder/bazel.sh
ARG BAZEL_VERSION=7.3.2
ARG DOCKER_VERSION=5:24.0.7-1~ubuntu.22.04~jammy

RUN \
    # Update package lists
    apt-get update && \
    apt-get -y install \
        python3 \
        python3-pkg-resources \
        software-properties-common \
        unzip \
        curl \
        ca-certificates \
        git && \

    # Install Java 11 (better compatibility with modern Bazel versions)
    apt-get -y install openjdk-11-jdk-headless && \

    # Install Bazel from binary release (https://github.com/bazelbuild/bazel/releases)
    # This is more reliable than the deprecated APT repository
    mkdir -p /tmp/bazel-install && \
    curl -L https://github.com/bazelbuild/bazel/releases/download/${BAZEL_VERSION}/bazel-${BAZEL_VERSION}-linux-x86_64 \
        -o /tmp/bazel-install/bazel-installer.sh && \
    chmod +x /tmp/bazel-install/bazel-installer.sh && \
    /tmp/bazel-install/bazel-installer.sh --prefix=/usr/local && \
    rm -rf /tmp/bazel-install && \

    # Install Docker (https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/#set-up-the-repository)
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | \
        gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
        https://download.docker.com/linux/ubuntu \
        $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list && \
    apt-get -y update && \
    apt-get install -y docker-ce=${DOCKER_VERSION} docker-ce-cli=${DOCKER_VERSION} && \

    # Cleanup
    rm -rf \
       /var/cache/debconf/* \
       /var/lib/apt/lists/* \
       /var/log/* \
       /tmp/* \
       /var/tmp/* && \

    # Install bazel.sh wrapper
    mv /builder/bazel.sh /usr/bin/bazel && \
    chmod +x /usr/bin/bazel && \

    # Test Bazel installation
    bazel version

# Create workspace directory (will be mounted by Cloud Build)
RUN mkdir -p /workspace

# Note: .bazelrc is no longer created statically here. Instead, bazel.sh
# sets up workspace configuration dynamically at runtime to handle permission
# and isolation issues across different Cloud Build execution contexts.

ENTRYPOINT ["bazel"]
