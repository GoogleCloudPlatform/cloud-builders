FROM golang:stretch AS build-env
WORKDIR /go-src
COPY go.mod /
COPY go.sum /
RUN go mod download
ADD . /go-src
RUN go build -o /gke-deploy

FROM google/cloud-sdk:alpine
RUN gcloud -q components install kubectl
RUN gcloud -q components install gsutil

COPY --from=build-env /gke-deploy /
COPY --from=build-env /gke-deploy /bin
COPY VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT [ "/gke-deploy" ]
