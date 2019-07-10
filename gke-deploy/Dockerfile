FROM golang:stretch AS build-env
ADD . /go-src
WORKDIR /go-src
RUN go build -o /gke-deploy

FROM google/cloud-sdk:alpine
RUN gcloud -q components install kubectl

COPY --from=build-env /gke-deploy /
COPY VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT [ "/gke-deploy" ]
