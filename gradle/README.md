# Tool builder: `gcr.io/cloud-builders/gradle`

This Cloud Build builder runs Gradle.

## Using an official [`gradle`](https://hub.docker.com/_/gradle) image

Because the official `gradle` image in Dockerhub specifies `USER gradle` and GCB
runs builds as `root`, the official `gradle` images are not currently directly
usable in GCB.

If you want to use a version of `gradle` that is not supported in this repo,
build an image using a Dockerfile like so:

```
FROM gradle:[VERSION]
ENV GRADLE_USER_HOME "/root/.gradle/"
USER root
ENTRYPOINT ["gradle"]
```

## Building this builder

To build this builder, run the following command in this directory.

    $ gcloud builds submit
