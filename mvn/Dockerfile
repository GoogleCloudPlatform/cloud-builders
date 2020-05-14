ARG MAVEN_VERSION=latest
FROM maven:${MAVEN_VERSION}
COPY deprecation.sh /bin
ENTRYPOINT ["/bin/deprecation.sh"]
