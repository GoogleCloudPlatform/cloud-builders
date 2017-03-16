load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

docker_build(
    name = "target",
    base = "@docker_debian//:wheezy",
    entrypoint = [
        "echo",
        "foo",
    ],
)
