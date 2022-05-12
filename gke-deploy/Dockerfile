FROM golang:stretch AS build-env
WORKDIR /go-src
COPY go.mod /
COPY go.sum /
RUN go mod download
ADD . /go-src
RUN go build -o /gke-deploy

FROM gcr.io/google.com/cloudsdktool/cloud-sdk:alpine
RUN gcloud -q components install kubectl
RUN gcloud -q components install gsutil
RUN gcloud -q components install kustomize
RUN gcloud -q components install local-extract
RUN apk -q --no-cache add gettext

COPY --from=build-env /gke-deploy /
COPY --from=build-env /gke-deploy /bin
COPY VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT [ "/gke-deploy" ]
