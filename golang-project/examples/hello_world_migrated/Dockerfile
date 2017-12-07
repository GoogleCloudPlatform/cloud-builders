FROM alpine
ENV PATH=/golang_project_bin:$PATH
COPY gopath/bin/hello /golang_project_bin/
ENTRYPOINT ["/golang_project_bin/hello"]
