FROM gcr.io/google_appengine/debian8

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
# + 1.0.4, which supporst .csproj projects for .NET Core 1.0.5 and 1.1.2.
RUN mkdir -p /usr/share/dotnet && \
    curl -sL https://download.microsoft.com/download/E/7/8/E782433E-7737-4E6C-BFBF-290A0A81C3D7/dotnet-dev-debian-x64.1.0.4.tar.gz | tar -xz -C /usr/share/dotnet/ && \
    curl -sL https://go.microsoft.com/fwlink/?LinkID=836302 | tar -xz -C /usr/share/dotnet && \
    ln -s /usr/share/dotnet/dotnet /usr/bin/dotnet

# Make sure that the caches are expanded.
ADD ./configs /configs
RUN cd /configs/config-1.0.0/ && dotnet new && \
    cd /configs/config-1.0.4/ && dotnet new && \
    rm -rf /configs

ENTRYPOINT [ "/usr/bin/dotnet" ]
