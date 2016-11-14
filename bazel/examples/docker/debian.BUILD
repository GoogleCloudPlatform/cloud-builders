load("@bazel_tools//tools/build_defs/docker:docker.bzl", "docker_build")

# Extract .xz files
genrule(
    name = "wheezy_tar",
    srcs = ["tianon-docker-brew-debian-e9bafb1/wheezy/rootfs.tar.xz"],
    outs = ["wheezy_tar.tar"],
    cmd = "cat $< | xzcat >$@",
)

docker_build(
    name = "wheezy",
    tars = [":wheezy_tar"],
    visibility = ["//visibility:public"],
)
