VERSION = 0.48.0
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
ifeq ($(OS),Windows_NT)
GOPATH_ROOT:=$(shell cygpath ${GOPATH})
else
GOPATH_ROOT:=${GOPATH}
endif

export GO111MODULE=on

.PHONY: all
all: clean testconvention test build rpm deb

.PHONY: test
test:
	go test $(TESTFLAGS) ./...
	./test.bash

.PHONY: lint
lint:
	golangci-lint run

.PHONY: testconvention
testconvention:
	prove -r t/
	@go generate ./... && git diff --exit-code -- . ':(exclude)go.*' || (echo 'please `go generate ./...` and commit them' && false)

.PHONY: build
build:
	mkdir -p build
	for i in $(filter-out check-windows-%, $(wildcard check-*)); do \
	  CGO_ENABLED=0 go build -ldflags "-s -w" -o build/$$i \
	  `pwd | sed -e "s|${GOPATH_ROOT}/src/||"`/$$i; \
	done

# This recipe will be called from some other recipes.
# Though these other recipes expects this to build multi platform binaries
# using GOOS, GOARCH or both, "make" itself don't recognize switching of
# environment variables, this recipe will just say "mackerel-check is up to date."
# We need to force rebuild "mackerel-check" if GOOS or GOARCH are passed.
build/mackerel-check: $(patsubst %,depends_on,$(GOOS)$(GOARCH))
	mkdir -p build
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.gitcommit=$(CURRENT_REVISION)" \
	  -o build/mackerel-check

.PHONY: depends_on
depends_on:
	@:

.PHONY: rpm
rpm: rpm-v2

.PHONY: rpm-v2
rpm-v2: rpm-v2-x86 rpm-v2-arm

.PHONY: rpm-v2-x86
rpm-v2-x86:
	make build/mackerel-check GOOS=linux GOARCH=amd64
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --define "dist .el7.centos" \
	  --target x86_64 -bb packaging/rpm/mackerel-check-plugins-v2.spec
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --define "dist .amzn2" \
	  --target x86_64 -bb packaging/rpm/mackerel-check-plugins-v2.spec

.PHONY: rpm-v2-arm
rpm-v2-arm:
	make build/mackerel-check GOOS=linux GOARCH=arm64
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch aarch64" --define "dist .el7.centos" \
	  --target aarch64 -bb packaging/rpm/mackerel-check-plugins-v2.spec
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch aarch64" --define "dist .amzn2" \
	  --target aarch64 -bb packaging/rpm/mackerel-check-plugins-v2.spec

.PHONY: deb
deb: deb-v2-x86 deb-v2-arm

.PHONY: deb-v2-x86
deb-v2-x86:
	make build/mackerel-check GOOS=linux GOARCH=amd64
	cp build/mackerel-check packaging/deb-v2/debian/
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us

.PHONY: deb-v2-arm
deb-v2-arm:
	make build/mackerel-check GOOS=linux GOARCH=arm64
	cp build/mackerel-check packaging/deb-v2/debian/
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us -aarm64

.PHONY: clean
clean:
	if [ -d build ]; then \
	  rm -f build/check-*; \
	  rmdir build; \
	fi
	go clean

.PHONY: update
update:
	go get -u ./...
	go mod tidy
