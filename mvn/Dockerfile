ARG MAVEN_VERSION=latest
FROM maven:${MAVEN_VERSION}
# Upstream maven images for JDK >= 16 are Oracle Linux (ol) based
# and others are based on Debian base image
RUN true && \
    v=`awk </etc/os-release -F = '$1 ~ /VERSION_ID/ {gsub(/"/,"",$2);print set $2}'` && \
    dist=`awk </etc/os-release -F = '$1 ~ /^ID_LIKE/ {gsub(/"/,"",$2);print set $2}'` && \
    if [ "$dist" = "debian" ] && [ "$v" -gt 9 ] ;  \
    then apt-get update -qqy && apt-get dist-upgrade -yq ; \
    elif  [ "$dist" = "fedora" ] ; \
    then microdnf clean all && microdnf update --refresh --nodocs && microdnf clean all ; \
    fi

ENTRYPOINT ["mvn"]
