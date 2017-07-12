ARG BASE_IMAGE=openjdk:8-jdk-alpine
FROM ${BASE_IMAGE}

RUN apk add --update --no-cache curl tar bash

ARG MAVEN_VERSION=3.5.0
ARG USER_HOME_DIR="/root"
ARG SHA=beb91419245395bd69a4a6edad5ca3ec1a8b64e41457672dc687c173a495f034
ARG BASE_URL=https://apache.osuosl.org/maven/maven-3/${MAVEN_VERSION}/binaries

RUN mkdir -p /usr/share/maven /usr/share/maven/ref \
  && curl -fsSL -o /tmp/apache-maven.tar.gz ${BASE_URL}/apache-maven-$MAVEN_VERSION-bin.tar.gz \
  && echo "${SHA}  /tmp/apache-maven.tar.gz" | sha256sum -c - \
  && tar -xzf /tmp/apache-maven.tar.gz -C /usr/share/maven --strip-components=1 \
  && rm -f /tmp/apache-maven.tar.gz \
  && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn

ENV M2_HOME /usr/share/maven

# Point maven central to google internal mirror
ADD settings.xml $USER_HOME_DIR/.m2/settings.xml

# transitively resolve all dependencies
ADD deps.txt /builder/deps.txt
ADD resolve-deps.sh /builder/resolve-deps.sh

RUN /builder/resolve-deps.sh

ENTRYPOINT ["mvn"]
