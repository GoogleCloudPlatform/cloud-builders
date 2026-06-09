FROM us-docker.pkg.dev/artifact-foundry-prod/docker-3p-trusted/python@sha256:f2ee145f3bc4e061f8dfe7e6ebd427a410121495a0bd26e7622136db060c59e0
RUN apt-get update -qqy && apt-get dist-upgrade -yq
COPY requirements.txt .
RUN /bin/sh -c set -eux; pip install --require-hashes -r requirements.txt --no-deps
ENTRYPOINT ["python3", "-m", "twine"]