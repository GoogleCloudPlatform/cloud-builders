FROM launcher.gcr.io/google/bazel

ADD bazel.sh /builder/bazel.sh

RUN \
    # This makes add-apt-repository available.
    apt-get update && \
    apt-get -y install \
        python \
        python-pkg-resources \
        software-properties-common \
        unzip && \

    # Install Git >2.0.1
    add-apt-repository ppa:git-core/ppa && \
    apt-get -y update && \
    apt-get -y install git && \

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
    apt-get install -y docker-ce=5:18.09.0~3-0~ubuntu-xenial unzip && \
    apt-get update && \

    orig=$(which bazel)                 && \
    mv $orig /builder/bazel             && \
    mv /builder/bazel.sh /usr/bin/bazel && \

    # Unpack bazel for future use.
    bazel version

# Store the Bazel outputs under /workspace so that the symlinks under bazel-bin (et al) are accessible
# to downstream build steps.
RUN mkdir -p /workspace
RUN echo 'startup --output_base=/workspace/.bazel' > ~/.bazelrc

ENTRYPOINT ["bazel"]
