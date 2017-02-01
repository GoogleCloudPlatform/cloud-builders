FROM gcr.io/cloud-builders/git

ADD write_version.bash /write_version.bash

RUN chmod +x /write_version.bash

ENTRYPOINT ["/write_version.bash"]
