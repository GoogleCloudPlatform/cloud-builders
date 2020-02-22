FROM golang:stretch AS build-env
ADD . /go-src
WORKDIR /go-src
RUN go build -o /gke-deploy

FROM google/cloud-sdk:alpine
RUN gcloud -q components install kubectl
RUN gcloud -q components install gsutil

COPY --from=build-env /gke-deploy /
COPY --from=build-env /gke-deploy /bin
COPY VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT [ "/gke-deploy" ]
