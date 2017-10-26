VERSION = 0.14.1
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
ifeq ($(OS),Windows_NT)
GOPATH_ROOT:=$(shell cygpath ${GOPATH})
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

check-release-deps:
	@have_error=0; \
	for command in cpanm hub ghch gobump; do \
	  if ! command -v $$command > /dev/null; then \
	    have_error=1; \
	    echo "\`$$command\` command is required for releasing"; \
	  fi; \
	done; \
	test $$have_error = 0

lint: devel-deps
	go vet ./...
	golint -set_exit_status ./...

testconvention:
	prove -r t/
	@go generate ./... && git diff --exit-code || (echo 'please `go generate ./...` and commit them' && false)

cover: devel-deps
	gotestcover -v -short -covermode=count -coverprofile=.profile.cov -parallelpackages=4 ./...

build: deps
	mkdir -p build
	for i in $(filter-out check-windows-%, $(wildcard check-*)); do \
	  go build -ldflags "-s -w" -o build/$$i \
	  `pwd | sed -e "s|${GOPATH_ROOT}/src/||"`/$$i; \
	done

build/mackerel-check: deps
	mkdir -p build
	go build -ldflags="-s -w -X main.gitcommit=$(CURRENT_REVISION)" \
	  -o build/mackerel-check

rpm: rpm-v1 rpm-v2

rpm-v1:
	make build GOOS=linux GOARCH=386
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" --define "buildarch noarch" -bb packaging/rpm/mackerel-check-plugins.spec
	make build GOOS=linux GOARCH=amd64
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" --define "buildarch x86_64" -bb packaging/rpm/mackerel-check-plugins.spec

rpm-v2:
	make build/mackerel-check GOOS=linux GOARCH=amd64
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --define "dist .el7.centos" \
	  -bb packaging/rpm/mackerel-check-plugins-v2.spec

deb: deb-v1 deb-v2

deb-v1:
	make build GOOS=linux GOARCH=386
	for i in `cat packaging/deb/debian/source/include-binaries`; do \
	  cp build/`basename $$i` packaging/deb/debian/; \
	done
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

deb-v2:
	make build/mackerel-check GOOS=linux GOARCH=amd64
	cp build/mackerel-check packaging/deb-v2/debian/
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us

release: check-release-deps
	(cd tool && cpanm -qn --installdeps .)
	perl tool/create-release-pullrequest

clean:
	if [ -d build ]; then \
	  rm -f build/check-*; \
	  rmdir build; \
	fi
	go clean

.PHONY: all test testconvention deps devel-deps lint cover build rpm rpm-v1 rpm-v2 deb deb-v1 deb-v2 clean release check-release-deps
