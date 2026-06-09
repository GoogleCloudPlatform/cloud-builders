FROM us-docker.pkg.dev/artifact-foundry-prod/docker-3p-trusted/golang@sha256:4746d26432a9117a5f58e95cb9f954ddf0de128e9d5816886514199316e4a2fb AS build-env
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
