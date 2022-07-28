FROM gcr.io/gcp-runtimes/ubuntu_20_0_4

RUN apt-get -y update && \
    apt-get -y install build-essential dpkg-dev fakeroot devscripts pbuilder cdbs && \
    # Clean up
    rm -rf /var/lib/apt/lists/*

