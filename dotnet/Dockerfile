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
# + 1.1.7, which supports .csproj projects for .NET Core 1.0.9 and 1.1.6.
# + 2.1.4, which supports .csproj projects for .NET Core 2.0.5.
RUN mkdir -p /usr/share/dotnet && \
    curl -sL https://go.microsoft.com/fwlink/?LinkID=836302 | tar -xz -C /usr/share/dotnet && \
    curl -sL https://download.microsoft.com/download/4/E/6/4E64A465-F02E-43AD-9A86-A08A223A82C3/dotnet-dev-debian-x64.1.1.7.tar.gz | tar -xz -C /usr/share/dotnet/ && \
    curl -sL https://download.microsoft.com/download/1/1/5/115B762D-2B41-4AF3-9A63-92D9680B9409/dotnet-sdk-2.1.4-linux-x64.tar.gz | tar -xz -C /usr/share/dotnet && \
    ln -s /usr/share/dotnet/dotnet /usr/bin/dotnet

# Make sure that the caches are expanded.
ADD ./configs /configs
RUN cd /configs/config-1.0.0/ && dotnet new && \
    cd /configs/config-1.1.7/ && dotnet new && \
    cd /configs/config-2.1.4/ && dotnet new && \
    rm -rf /configs

# Ensure that no additional cache expansions are performed from now on, see
# https://github.com/dotnet/cli/issues/6550 and
# https://github.com/dotnet/cli/pull/6633 for why this is necessary.
ENV DOTNET_SKIP_FIRST_TIME_EXPERIENCE=1

ENTRYPOINT [ "/usr/bin/dotnet" ]
