# Use the base App Engine Docker image, based on debian jessie.
FROM launcher.gcr.io/google/debian8

# https://yarnpkg.com/en/docs/install#linux-tab
RUN apt-get update && \
    apt-get install -y curl apt-transport-https git bzip2 build-essential && \

    # Install newer Node version.
    curl -sL https://deb.nodesource.com/setup_8.x | bash - && \
    apt-get install -y nodejs && \

    # Install Yarn
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update && \
    apt-get install -y yarn

ENTRYPOINT ["yarn"]
