ARG NODE_VERSION=latest
FROM node:${NODE_VERSION}
COPY ./deprecation.sh /bin
ENTRYPOINT ["/bin/deprecation.sh"]
