FROM maven:3.5.0-jdk-8

ADD poms /poms

# for each included pom file, resolve all of its dependencies
RUN for POM in $(find /poms/*.xml); do \
  mvn --batch-mode -f $POM -s /usr/share/maven/ref/settings-docker.xml org.apache.maven.plugins:maven-dependency-plugin:3.0.0:go-offline -DincludeParents ; \
  done; \
  rm -r /poms

