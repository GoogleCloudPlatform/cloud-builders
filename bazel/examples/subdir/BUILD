load(
   "@io_bazel_rules_docker//docker:docker.bzl",
   "docker_build",
)

docker_build(
    name = "target",
    base = "@base_image//image:image.tar",
    entrypoint = [
        "echo",
        "foo",
    ],
)
