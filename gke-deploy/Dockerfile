FROM golang:1.21 AS build-env
WORKDIR /go-src
COPY go.mod /
COPY go.sum /
RUN go mod download
ADD . /go-src
RUN CGO_ENABLED=0 go test ./...
RUN CGO_ENABLED=0 go build -o /gke-deploy

FROM gcr.io/google.com/cloudsdktool/cloud-sdk:alpine
RUN gcloud -q components install gke-gcloud-auth-plugin
RUN gcloud -q components install kubectl
RUN gcloud -q components install gsutil
RUN gcloud -q components install kustomize
RUN gcloud -q components install nomos
RUN gcloud -q components install local-extract
RUN apk update && apk upgrade --available --no-cache
RUN apk -q --no-cache add gettext
RUN apk -q --no-cache add yq


COPY --from=build-env /gke-deploy /
COPY --from=build-env /gke-deploy /bin
COPY VENDOR-LICENSE /
COPY LICENSE /
ENTRYPOINT [ "/gke-deploy" ]
