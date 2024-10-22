.PHONY: default check test build image

IMAGE_NAME := sablierapp/mimic

build:
	CGO_ENABLED=0 go build -a --trimpath --installsuffix cgo --ldflags="-s" -o mimic

check:
	golangci-lint run

image:
	docker build -t $(IMAGE_NAME) .