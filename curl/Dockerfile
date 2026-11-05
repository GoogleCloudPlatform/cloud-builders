FROM gcr.io/gcp-runtimes/ubuntu_20_0_4

ARG USE_PUBLIC_DEPS="true"

# For internal builds, we:
# 1. Backup standard sources temporarily
# 2. Configure Apt to use internal mirrors
# 3. Update and upgrade packages
# 4. Restore original sources
# 5. Clean up all temporary files and caches
# Run all apt operations in a single RUN instruction to keep the image size minimal.
RUN if [ "${USE_PUBLIC_DEPS}" = "true" ]; then \
        # Public build: use default public Ubuntu mirrors
        apt-get update && \
        apt-get dist-upgrade -y ; \
    else \
        # Internal build: use internal Ubuntu mirrors
        if [ -f /etc/apt/sources.list ]; then mv /etc/apt/sources.list /etc/apt/sources.list.bak; fi && \
        if [ -d /etc/apt/sources.list.d ]; then mv /etc/apt/sources.list.d /etc/apt/sources.list.d.bak; fi && \
        mkdir -p /etc/apt/sources.list.d && \
        . /etc/os-release && \
        echo "deb https://us-apt.pkg.dev/remote/artifact-foundry-prod/ubuntu-3p-remote-${VERSION_CODENAME} ${VERSION_CODENAME} main universe" > /etc/apt/sources.list.d/artifact-registry.list && \
        echo "deb https://us-apt.pkg.dev/remote/artifact-foundry-prod/ubuntu-3p-remote-${VERSION_CODENAME}-updates ${VERSION_CODENAME}-updates main universe" >> /etc/apt/sources.list.d/artifact-registry.list && \
        echo "deb https://us-apt.pkg.dev/remote/artifact-foundry-prod/ubuntu-3p-remote-${VERSION_CODENAME}-security ${VERSION_CODENAME}-security main universe" >> /etc/apt/sources.list.d/artifact-registry.list && \
        apt-get update && \
        apt-get dist-upgrade -y && \
        rm -rf /etc/apt/sources.list.d && \
        if [ -d /etc/apt/sources.list.d.bak ]; then mv /etc/apt/sources.list.d.bak /etc/apt/sources.list.d; else mkdir -p /etc/apt/sources.list.d; fi && \
        if [ -f /etc/apt/sources.list.bak ]; then mv /etc/apt/sources.list.bak /etc/apt/sources.list; fi ; \
    fi && \
    \
    # 4. Clean up all temporary files and caches
    apt-get clean && \
    rm -rf \
       /var/cache/debconf/* \
       /var/lib/apt/lists/* \
       /var/log/* \
       /tmp/* \
       /var/tmp/*
       
COPY notice.sh /usr/bin
ENTRYPOINT ["/usr/bin/notice.sh"]