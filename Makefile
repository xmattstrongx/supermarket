deploy: test build run

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell GCOMMIT=`git rev-parse --short HEAD`; if [ -n "`git status . --porcelain`" ]; then echo "$$GCOMMIT-dirty"; else echo $$GCOMMIT; fi)

ifeq "$(BUILD_NUMBER)" ""
    VERSION_SUFFIX   ?= $(GIT_COMMIT)
else
    VERSION_SUFFIX   ?= $(GIT_COMMIT)-$(CIRCLE_BUILD_NUM)
endif

VERSION := $(shell cat VERSION)-$(VERSION_SUFFIX)

## Devflow targets
test:
	go test -v ./...

build:
	echo building image
	docker build . -t supermarket:$(VERSION) || exit 1;

run:
	docker run --rm -it -p 8080:8080 supermarket:$(VERSION)

install:
	go install github.com/xmattstrongx/supermarket

## CI targets
ci-test: test

ci-docker-login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
	
ci-build:
	echo building image
	docker build . \
		--label version=$(VERSION) \
		-t xmattstrongx/supermarket:$(VERSION) \
		|| exit 1;

ci-push:
	echo pushing image
	docker push xmattstrongx/supermarket:$(VERSION) \
		|| exit 1;
	
	if [ $(GIT_BRANCH) = "master" ]; then \
            echo pushing new tag xmattstrongx/supermarket:$(VERSION) as latest; \
            docker tag xmattstrongx/supermarket:$(VERSION) xmattstrongx/supermarket:latest || exit 1; \
            docker push xmattstrongx/supermarket:latest || exit 1; \
        fi; \