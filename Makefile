TARGET_OSARCH="linux/386"

all: clean test build rpm deb

test: lint
	go test $(TESTFLAGS) ./...

deps:
	go get -d -v -t ./...

devel-deps: deps
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/vet
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

LINT_RET = .golint.txt
lint: devel-deps
	go vet ./...
	rm -f $(LINT_RET)
	golint ./... | tee -a $(LINT_RET)
	test ! -s $(LINT_RET)

cover: devel-deps
	tool/cover.sh

build: deps
	mkdir -p build
	for i in check-*; do \
	  if [ $$i = check-load ]; then \
	    echo ""; \
	  else \
	  gox -osarch=$(TARGET_OSARCH) -output build/$$i \
	    github.com/mackerelio/go-check-plugins/$$i; \
		fi; \
	done

rpm: build
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-checks.spec

deb: build
	cp build/check-* packaging/deb/debian/
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

clean:
	if [ -d build ]; then \
	  rm -f build/check-*; \
	  rmdir build; \
	fi
	go clean

.PHONY: all test deps devel-deps lint cover build clean rpm deb
