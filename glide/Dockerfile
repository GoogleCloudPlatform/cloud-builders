FROM gcr.io/test-argo/go:wheezy

RUN GOPATH=/go && curl https://glide.sh/get | sh

COPY glide.bash /bin/glide.bash

ENTRYPOINT ["/bin/glide.bash"]
