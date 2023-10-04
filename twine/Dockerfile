FROM python:3.10.7-slim-bullseye
RUN apt-get update -qqy && apt-get dist-upgrade -yq 
RUN /bin/sh -c set -eux; pip install twine==4.0.1
RUN /bin/sh -c set -eux; pip install keyrings.google-artifactregistry-auth==1.1.1
ENTRYPOINT ["python3", "-m", "twine"]
