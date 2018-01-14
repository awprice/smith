http_archive(
    name = "bazel_gazelle",
    sha256 = "e3dadf036c769d1f40603b86ae1f0f90d11837116022d9b06e4cd88cae786676",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.8/bazel-gazelle-0.8.tar.gz",
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "4d8d6244320dd751590f9100cf39fd7a4b75cd901e1f3ffdfd6f048328883695",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.9.0/rules_go-0.9.0.tar.gz",
)

git_repository(
    name = "io_bazel_rules_docker",
    commit = "3caf72f166f8b6b0e529442477a74871ad4d35e9",
    remote = "https://github.com/bazelbuild/rules_docker.git",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
load(
    "@io_bazel_rules_docker//go:image.bzl",
    go_image_repositories = "repositories",
)
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_rules_dependencies()

go_register_toolchains()

go_image_repositories()

gazelle_dependencies()
