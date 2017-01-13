FROM golang:alpine

# Install VCS tools to support "go get" commands.
RUN apk add --update --no-cache git mercurial subversion

# We blank out the GOPATH because the base image sets it, and
# if the user of this build step does *not* set it, we want to
# explore other options for workspace derivation.
ENV GOPATH=
# Only allow static linking, so that binaries created using
# this alpine base can still be run outside alpine environments.
ENV CGO_ENABLED=0

RUN mkdir /builder

COPY go_workspace.go prepare_workspace.inc go.sh /builder/

RUN go build -o /builder/go_workspace /builder/go_workspace.go

ENTRYPOINT ["/builder/go.sh"]
