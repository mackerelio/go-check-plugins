TARGET_OSARCH="linux/amd64"
CURRENT_VERSION = $(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')
ifeq ($(OS),Windows_NT)
GOPATH_ROOT:=$(shell cygpath ${GOPATH})
TARGET_OSARCH="windows/amd64"
else
GOPATH_ROOT:=${GOPATH}
endif

check-variables:
	echo "CURRENT_VERSION: ${CURRENT_VERSION}"

all: clean test build rpm deb

test: lint
	go test $(TESTFLAGS) ./...

deps:
	go get -d -v -t ./...

devel-deps: deps
	go get github.com/golang/lint/golint
	go get github.com/pierrre/gotestcover
	go get github.com/mattn/goveralls

LINT_RET = .golint.txt
lint: devel-deps
	go vet ./...
	rm -f $(LINT_RET)
	golint ./... | tee -a $(LINT_RET)
	test ! -s $(LINT_RET)

cover: devel-deps
	gotestcover -v -short -covermode=count -coverprofile=.profile.cov -parallelpackages=4 ./...

build: deps
	mkdir -p build
	for i in check-*; do \
	  gox -ldflags "-s -w" \
	    -osarch=$(TARGET_OSARCH) -output build/$$i \
	    `pwd | sed -e "s|${GOPATH_ROOT}/src/||"`/$$i; \
	done

rpm: build
	make build TARGET_OSARCH="linux/386"
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${CURRENT_VERSION}" --define "buildarch noarch" -bb packaging/rpm/mackerel-check-plugins.spec
	make build TARGET_OSARCH="linux/amd64"
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${CURRENT_VERSION}" --define "buildarch x86_64" -bb packaging/rpm/mackerel-check-plugins.spec

deb: deps
	TARGET_OSARCH="linux/386" make build
	cp build/check-* packaging/deb/debian/
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

clean:
	if [ -d build ]; then \
	  rm -f build/check-*; \
	  rmdir build; \
	fi
	go clean

.PHONY: all test deps devel-deps lint cover build rpm deb clean release
