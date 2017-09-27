ARG BASE_IMAGE=gcr.io/cloud-builders/java/javac:8
FROM ${BASE_IMAGE}

ARG GRADLE_VERSION=4.2
ARG USER_HOME_DIR="/root"
ARG SHA=515dd63d32e55a9c05667809c5e40a947529de3054444ad274b3b75af5582eae
ARG BASE_URL=https://services.gradle.org/distributions

ENV GRADLE_HOME "/usr/share/gradle-${GRADLE_VERSION}"
ENV GRADLE_USER_HOME "${USER_HOME_DIR}/.gradle/"

RUN apt-get update -qqy && apt-get install -qqy curl \
  && mkdir -p /usr/share "${GRADLE_USER_HOME}" \
  && curl -fsSL -o "gradle-${GRADLE_VERSION}-bin.zip" "${BASE_URL}/gradle-${GRADLE_VERSION}-bin.zip" \
  && echo "${SHA}  gradle-${GRADLE_VERSION}-bin.zip" | sha256sum -c - \
  && unzip -qq "gradle-${GRADLE_VERSION}-bin.zip" \
  && rm -f "gradle-${GRADLE_VERSION}-bin.zip" \
  && mv "gradle-${GRADLE_VERSION}" /usr/share \
  && ln -s "${GRADLE_HOME}/bin/gradle" /usr/bin/gradle \
  && apt-get remove -qqy --purge curl \
  && rm /var/lib/apt/lists/*_*

ADD gradle.properties "${GRADLE_USER_HOME}"

ENTRYPOINT ["/usr/bin/gradle"]

