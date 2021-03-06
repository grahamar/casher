.PHONY: all install-dependencies build

APPLICATION_ID=casher
ORG_PATH=github.com/grahamar
REPO_PATH=$(ORG_PATH)/$(APPLICATION_ID)
export PATH := $(PWD)/bin:$(PATH)

VERSION ?= $(shell git describe --tags --dirty --always | sed -e 's/^v//g')

BIN_NAME=dist/${APPLICATION_ID}
DOCKER_IMAGE=grahamar/$(APPLICATION_ID):$(VERSION)

export GOBIN=$(PWD)/bin

all: install-dependencies build

install-dependencies:
	@dep ensure -vendor-only

.PHONY: release-binary
release-binary: install-dependencies
	@go build -o /go/bin/$(APPLICATION_ID) -v $(REPO_PATH)

build:
	go build -o ${BIN_NAME}
	@echo "You can now use ./${BIN_NAME}"

docker-build:
	docker build -t $(DOCKER_IMAGE) --build-arg VERSION=$(VERSION) .

deploy: docker-build
	docker push $(DOCKER_IMAGE)
