#!/usr/bin/env bash

mkdir -p ./tmp/protoc
pushd ./tmp/protoc

CPU_ARCH=""
case $(uname -m) in
    "x86_64") CPU_ARCH="x86_64" ;;
    "aarch64") CPU_ARCH="aarch_64" ;;
esac

PROTOC_VERSION="3.12.4"
PROTOC_ARCHIVE="protoc-${PROTOC_VERSION}-linux-${CPU_ARCH}.zip"

curl -fsSLO "https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ARCHIVE}"
unzip "$PROTOC_ARCHIVE"
sudo mv bin/protoc /usr/local/bin
popd

rm -rf ./tmp
