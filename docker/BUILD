load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_image",
)

container_image(
    name = "docker-17.05.0",
    base = "@base_image//image:image.tar",
    directory = "/bin",
    entrypoint = ["/bin/docker"],
    files = ["@docker_static_17_05_0//:docker_static"],
    mode = "0777",
)

container_image(
    name = "docker-17.06.1",
    base = "@base_image//image:image.tar",
    directory = "/bin",
    entrypoint = ["/bin/docker"],
    files = ["@docker_static_17_06_1//:docker_static"],
    mode = "0777",
)

# TODO(jasonhall): Use container_bundle which lets us build and tag multiple
# images in one rule
# (https://github.com/bazelbuild/rules_docker#container_bundle-1). This depends
# on us having some way to pass $PROJECT_ID into the image tag.
