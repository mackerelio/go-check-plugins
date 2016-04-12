TARGET_OSARCH="linux/amd64"

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
	    github.com/mackerelio/go-check-plugins/$$i; \
	done

rpm: build
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-check-plugins.spec

deb: build
	cp build/check-* packaging/deb/debian/
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

clean:
	if [ -d build ]; then \
	  rm -f build/check-*; \
	  rmdir build; \
	fi
	go clean

release:
	tool/releng

.PHONY: all test deps devel-deps lint cover build rpm deb clean release
