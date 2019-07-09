FROM golang:stretch AS build-env
RUN go get github.com/gobuffalo/packr/v2/packr2
ADD . /go/src/github.com/GoogleCloudPlatform/cloud-builders/gke-deploy
RUN go get -d github.com/GoogleCloudPlatform/cloud-builders/gke-deploy
RUN cd /go/src/github.com/GoogleCloudPlatform/cloud-builders/gke-deploy && packr2
RUN go build -o /gke-deploy github.com/GoogleCloudPlatform/cloud-builders/gke-deploy

FROM google/cloud-sdk:alpine
RUN gcloud -q components install kubectl

COPY --from=build-env /gke-deploy /
COPY VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT [ "/gke-deploy" ]
