IMAGE_NAME ?= kiada-go
IMAGE_TAG  ?= 1.0
REGISTRY   ?= docker.io/$(shell id -un)

.PHONY: build run test docker-build docker-push

build:
	go build -o kiada-go .

run:
	go run . 

test:
	go test ./...

docker-build:
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .

docker-push: docker-build
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

clean:
	rm -f kiada-go
