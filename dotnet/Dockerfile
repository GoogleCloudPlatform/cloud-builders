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
# + 1.1.4, which supporst .csproj projects for .NET Core 1.0.7 and 1.1.4.
# + 2.0.0, which supports .csproj projects for .NET Core 2.0.0.
RUN mkdir -p /usr/share/dotnet && \
    curl -sL https://go.microsoft.com/fwlink/?LinkID=836302 | tar -xz -C /usr/share/dotnet && \
    curl -sL https://download.microsoft.com/download/F/4/F/F4FCB6EC-5F05-4DF8-822C-FF013DF1B17F/dotnet-dev-debian-x64.1.1.4.tar.gz | tar -xz -C /usr/share/dotnet/ && \
    curl -sL https://download.microsoft.com/download/1/B/4/1B4DE605-8378-47A5-B01B-2C79D6C55519/dotnet-sdk-2.0.0-linux-x64.tar.gz | tar -xz -C /usr/share/dotnet && \
    ln -s /usr/share/dotnet/dotnet /usr/bin/dotnet

# Make sure that the caches are expanded.
ADD ./configs /configs
RUN cd /configs/config-1.0.0/ && dotnet new && \
    cd /configs/config-1.1.4/ && dotnet new && \
    cd /configs/config-2.0.0/ && dotnet new && \
    rm -rf /configs

# Ensure that no additional cache expansions are performed from now on, see
# https://github.com/dotnet/cli/issues/6550 and
# https://github.com/dotnet/cli/pull/6633 for why this is necessary.
ENV DOTNET_SKIP_FIRST_TIME_EXPERIENCE=1

ENTRYPOINT [ "/usr/bin/dotnet" ]
