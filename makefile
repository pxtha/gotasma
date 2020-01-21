PROJECT_NAME=gotasma
BUILD_VERSION=$(shell cat VERSION)
DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

REALIZE_VERSION=2.0.2
REALIZE_IMAGE=realize:$(REALIZE_VERSION)

.SILENT:

all: mod_tidy fmt vet install

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin .

install:
	$(GO_BUILD_ENV) go install

vet:
	$(GO_BUILD_ENV) go vet $(GO_FILES)

fmt:
	$(GO_BUILD_ENV) go fmt $(GO_FILES)

mod_tidy:
	$(GO_BUILD_ENV) go mod tidy

compose_dev: realize
	cd deployment/dev && REALIZE_VERSION=$(REALIZE_VERSION) docker-compose up

realize:
	cd deployment/dev; \
	docker build -t $(REALIZE_IMAGE) .;

docker_run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)

compose_prod: docker
	cd deployment/docker && BUILD_VERSION=$(BUILD_VERSION) docker-compose up

docker_prebuild: vet build
	mkdir -p deployment/docker/configs
	mv $(PROJECT_NAME)-$(BUILD_VERSION).bin deployment/docker/$(PROJECT_NAME).bin; \
	cp -R configs deployment/docker/;

docker_build:
	cd deployment/docker; \
	docker build -t $(DOCKER_IMAGE) .;

docker_postbuild:
	cd deployment/docker; \
	rm -rf $(PROJECT_NAME).bin 2> /dev/null;\
	rm -rf configs 2> /dev/null;

docker: docker_prebuild docker_build docker_postbuild