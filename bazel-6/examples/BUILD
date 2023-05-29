load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "target_default_library",
    srcs = ["main.go"],
    importpath = "github.com/GoogleCloudPlatform/cloud-builders/bazel/examples",
    visibility = ["//visibility:private"],
    deps = [],
)

go_binary(
    name = "target-bin",
    embed = [":target_default_library"],
    pure = "on",
    visibility = ["//visibility:public"],
)

load("@io_bazel_rules_docker//go:image.bzl", "go_image")

go_image(
    name = "target",
    binary = ":target-bin",
    visibility = ["//visibility:public"],
)

go_library(
    name = "checkargs_default_library",
    srcs = ["main_checkargs.go"],
    importpath = "github.com/GoogleCloudPlatform/cloud-builders/bazel/examples",
    visibility = ["//visibility:private"],
    deps = [],
)

go_binary(
    name = "checkargs",
    embed = [":checkargs_default_library"],
    pure = "on",
    visibility = ["//visibility:public"],
)
