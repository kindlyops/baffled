#!/bin/bash

set -e

sha256=7f8e26b1759ab52ff8aef192e9d92cf5a216746721ca420f4b4a32a25e3d4b91
version=0.13.0
file=bazel-${version}-without-jdk-installer-linux-x86_64.sh
curl -sLO https://github.com/bazelbuild/bazel/releases/download/${version}/${file}
printf "${sha256}\t*${file}" | sha256sum --strict -c -
chmod +x ${file}
bash ${file} --user