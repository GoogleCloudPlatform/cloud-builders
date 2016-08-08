# Dockerizer

This is a builder to simply invoke `docker build`.

You must supply the name of the tag as the only argument.

To override the location of the `Dockerfile` to build, set the environment
variable `DOCKERFILE`.

The version of Docker that is used by this builder is `1.9.1`.
