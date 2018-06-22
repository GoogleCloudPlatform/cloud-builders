FROM launcher.gcr.io/google/debian8

# Install the dependencies for running the .NET Core SDK.
RUN apt-get update \
    && apt-get install -y \
        ca-certificates \
        curl \
        libc6 \
        libcurl3 \
        libgcc1 \
        libicu52 \
        liblttng-ust0 \
        libssl1.0.0 \
        libstdc++6 \
        libtinfo5 \
        libunwind8 \
        libuuid1 \
        zlib1g \
    && apt-get clean

# Install the following SDKs:
# + 1.0.3 with SDK Preview 2 build 3156, to support .NET Core 1.0.5 with project.json solutions.
# + 1.1.7, which supports .csproj projects for .NET Core 1.0.11 and 1.1.8.
# + 2.1.300 which supports .NET Core 2.1.0.
RUN mkdir -p /usr/share/dotnet && \
    curl -sL https://download.microsoft.com/download/0/3/0/030449F5-F093-44A6-9889-E19B50A59777/sdk/dotnet-dev-debian-x64.1.0.0-preview2-003156.tar.gz | tar -xz -C /usr/share/dotnet && \
    curl -sL https://download.microsoft.com/download/4/0/2/4022CFC7-5061-4762-B9BA-48B35632572D/dotnet-dev-debian-x64.1.1.9.tar.gz | tar -xz -C /usr/share/dotnet/ && \
    curl -sL https://download.microsoft.com/download/D/0/4/D04C5489-278D-4C11-9BD3-6128472A7626/dotnet-sdk-2.1.301-linux-x64.tar.gz | tar -xz -C /usr/share/dotnet && \
    ln -s /usr/share/dotnet/dotnet /usr/bin/dotnet

# Make sure that the caches are expanded.
ADD ./configs /configs
RUN cd /configs/config-1.0.0/ && dotnet new && \
    cd /configs/config-1.1.9/ && dotnet new && \
    cd /configs/config-2.1.301/ && dotnet new && \
    rm -rf /configs

# Ensure that no additional cache expansions are performed from now on, see
# https://github.com/dotnet/cli/issues/6550 and
# https://github.com/dotnet/cli/pull/6633 for why this is necessary.
ENV DOTNET_SKIP_FIRST_TIME_EXPERIENCE=1

ENTRYPOINT [ "/usr/bin/dotnet" ]
