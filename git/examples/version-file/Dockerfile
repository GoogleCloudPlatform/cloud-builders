FROM gcr.io/cloud-builders/git

ADD write_version.ash /write_version.ash

RUN chmod +x /write_version.ash

ENTRYPOINT ["/write_version.ash"]
