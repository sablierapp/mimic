.PHONY: default check test build image

IMAGE_NAME := sablierapp/mimic

build:
	goreleaser build --single-target --clean --snapshot

check:
	golangci-lint run

image:
	docker build -t $(IMAGE_NAME) .