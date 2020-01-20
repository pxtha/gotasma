PROJECT_NAME=Gotasma
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