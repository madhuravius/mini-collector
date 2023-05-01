GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
SHELL=/bin/bash

.PHONY: deps
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go mod download

.PHONY: build
build: $(GOFILES_NOVENDOR)
	go list ./...  | grep cmd | xargs -P $$(nproc) -n 1 -- go build

protobufs:
	protoc -I protobufs/ protobufs/api.proto \
		--go-grpc_out=protobufs \
		--plugin=protoc-gen-custom="$PWD/.codegen/codegen" \
		--go_out=protobufs

.PHONY: gofiles
src: $(GOFILES_NOVENDOR) fmt
	@true

.PHONY: unit
unit: $(GOFILES_NOVENDOR)
	go test $$(go list ./...)
	go vet $$(go list ./...)

.PHONY: test
test: unit
	@true

.PHONY: fmt
fmt:
	gofmt -l -w ${GOFILES_NOVENDOR}

.DEFAULT_GOAL := test
