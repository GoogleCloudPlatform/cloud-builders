# Go boilerplate
load("@io_bazel_rules_go//go:def.bzl", "go_binary")
load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_build",
)

go_binary(
    name = "hello",
    srcs = ["main.go"],
)

docker_build(
    name = "target",
    base = "@base_image//image:image.tar",
    entrypoint = ["/hello"],
    files = [":hello"],
)
