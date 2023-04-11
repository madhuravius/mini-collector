#!/usr/bin/env bash

PYENV_VERSION=3.6

pip3 install --user -r .codegen/requirements.txt
mkdir -p /tmp/protoc
pushd /tmp/protoc

PROTOC_VERSION="3.12.4"
PROTOC_ARCHIVE="protoc-${PROTOC_VERSION}-linux-x86_64.zip"

curl -fsSLO "https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ARCHIVE}"
unzip "$PROTOC_ARCHIVE"
rm "$PROTOC_ARCHIVE"
