FROM python:3.10.7-slim-bullseye
RUN apt-get update -qqy && apt-get dist-upgrade -yq
COPY requirements.txt .
RUN /bin/sh -c set -eux; pip install --require-hashes -r requirements.txt --no-deps
ENTRYPOINT ["python3", "-m", "twine"]