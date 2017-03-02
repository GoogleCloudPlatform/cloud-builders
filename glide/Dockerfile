FROM gcr.io/test-argo/go:wheezy

ENV GOPATH=/go
RUN curl https://glide.sh/get | sh
ENV GOPATH=

COPY glide.bash /bin/glide.bash

ENTRYPOINT ["/bin/glide.bash"]
