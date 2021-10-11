.PHONY: all
default: build;

fmt:
	go fmt ./...

tests:
	go test ./...

build:
	goreleaser build --snapshot --skip-validate --rm-dist --single-target