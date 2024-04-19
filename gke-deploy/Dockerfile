FROM golang:1.21 AS build-env
WORKDIR /go-src
COPY go.mod /
COPY go.sum /
RUN go mod download
ADD . /go-src
RUN CGO_ENABLED=0 go test ./...
RUN CGO_ENABLED=0 go build -o /gke-deploy

FROM gcr.io/google.com/cloudsdktool/cloud-sdk:alpine

COPY --from=build-env /gke-deploy /bin
COPY VENDOR-LICENSE /
COPY LICENSE /

RUN gcloud -q components install \
    gke-gcloud-auth-plugin \
    kubectl \
    gsutil \
    kustomize \
    nomos \
    local-extract \
    && rm -rf /google-cloud-sdk/.install/.backup \
    && apk update && apk upgrade --available --no-cache \
    && apk -q --no-cache add \
        gettext \
        yq \
    && ln  /bin/gke-deploy /gke-deploy

ENTRYPOINT [ "/bin/gke-deploy" ]
