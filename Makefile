TARGET_OSARCH="linux/amd64"
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
ifeq ($(OS),Windows_NT)
GOPATH_ROOT:=$(shell cygpath ${GOPATH})
TARGET_OSARCH="windows/amd64"
else
GOPATH_ROOT:=${GOPATH}
endif

all: clean testconvention test build rpm deb

test: lint
	go test $(TESTFLAGS) ./...

deps:
	go get -d -v -t ./...

devel-deps: deps
	go get github.com/golang/lint/golint
	go get github.com/pierrre/gotestcover
	go get github.com/mattn/goveralls

lint: devel-deps
	go vet ./...
	golint -set_exit_status ./...

testconvention:
	prove -r t/
	@go generate ./... && git diff --exit-code || (echo 'please `go generate ./...` and commit them' && exit 1)

cover: devel-deps
	gotestcover -v -short -covermode=count -coverprofile=.profile.cov -parallelpackages=4 ./...

build: deps
	mkdir -p build
	for i in $(filter-out check-windows-%, $(wildcard check-*)); do \
	  gox -ldflags "-s -w" \
	    -osarch=$(TARGET_OSARCH) -output build/$$i \
	    `pwd | sed -e "s|${GOPATH_ROOT}/src/||"`/$$i; \
	done

build/mackerel-check: deps
	mkdir -p build
	go build -ldflags="-s -w -X main.gitcommit=$(CURRENT_REVISION)" \
	  -o build/mackerel-check

rpm: build
	make build TARGET_OSARCH="linux/386"
	rpmbuild --define "_sourcedir `pwd`"  --define "_version 0.10.2" --define "buildarch noarch" -bb packaging/rpm/mackerel-check-plugins.spec
	make build TARGET_OSARCH="linux/amd64"
	rpmbuild --define "_sourcedir `pwd`"  --define "_version 0.10.2" --define "buildarch x86_64" -bb packaging/rpm/mackerel-check-plugins.spec

deb: deps
	TARGET_OSARCH="linux/386" make build
	cp build/check-* packaging/deb/debian/
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

release:
	(cd tool && cpanm -qn --installdeps .)
	perl tool/create-release-pullrequest

clean:
	if [ -d build ]; then \
	  rm -f build/check-*; \
	  rmdir build; \
	fi
	go clean

.PHONY: all test testconvention deps devel-deps lint cover build rpm deb clean release
