PROJECT=user-service
IMAGENAME=$(PROJECT)
DOCKERHUB_REPOSITORY=tochkaksn/$(IMAGENAME)
VERSION?=$(shell git describe --tags --always)
TARGET_OS ?= $(shell go env GOOS)
GIT_BRANCH?=$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

BINARY?=app
ifeq ($(TARGET_OS),windows)
BINARY+=$(BINARY).exe
endif

DOCKER_COMPOSE_OPTIONS=-f ./service-test/docker-compose.yml

ifndef GOPATH
GOPATH=~/go
endif


LDF_FLAGS=-X main.version=$(VERSION)
ifneq ($(TARGET_OS),darwin)
LDF_FLAGS+= -extldflags "-static"
endif

BUILD_ARGS=--ldflags '$(LDF_FLAGS)'

.DEFAULT_GOAL := build

define tag_docker
  @if [ "$(GIT_BRANCH)" = "master" ]; then \
    docker tag $(IMAGENAME) $(1):latest; \
  fi
	@if [ "$(GIT_BRANCH)" = "develop" ]; then \
    docker tag $(IMAGENAME) $(1):unstable; \
  fi
  @if [ "$(GIT_BRANCH)" != "master" && "$(GIT_BRANCH)" != "develop" ]; then \
    docker tag $(IMAGENAME) $(1):$(GIT_BRANCH); \
  fi
endef

# TEST SECTION
.PHONY: test test-all

test:
	@go test ./src/...

test-all:
	@go test  -tags=integration ./src/...


# BUILD SECTION
.PHONY: clean get build build_in_docker-env

clean:
	@rm -f $(BINARY)

$(GOPATH)/bin/dep:
	@go get github.com/golang/dep/cmd/dep

get: $(GOPATH)/bin/dep
	@dep ensure

build: get clean $(BINARY)

$(BINARY):
	@GOOS=$(TARGET_OS) go build  $(BUILD_ARGS) -o $(BINARY) ./src


# DOCKER SECTION
.PHONY: docker-rebuild docker-build docker-publish docker-dockerhub-publish docker-bintray-publish docker-inspect

docker-rebuild: clean docker-build

docker-build: $(BINARY)
	@docker build -t $(IMAGENAME) --build-arg VERSION=$(VERSION) .

docker-publish: docker-dockerhub-publish

docker-dockerhub-publish:
	$(call tag_docker, $(DOCKERHUB_REPOSITORY))
	docker push $(DOCKERHUB_REPOSITORY)

docker-up:
	@docker-compose $(DOCKER_COMPOSE_OPTIONS) up -d
	@sleep 15

docker-down:
	@docker-compose $(DOCKER_COMPOSE_OPTIONS) down

docker-test:docker-up service-test docker-down

docker-inspect:
	@docker inspect -f '{{index .ContainerConfig.Labels "VERSION"}}' $(IMAGENAME)

.PHONY: service-test

service-test:
	@go test ./service-test -v
