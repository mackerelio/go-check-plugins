TARGET_OSARCH="linux/386"

all: clean test build

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
	  gox -osarch=$(TARGET_OSARCH) -output build/$$i \
	    github.com/mackerelio/go-check-plugins/$$i; \
	done

clean:
	if [ -d build ]; then \
	  rm -f build/check-*; \
	  rmdir build; \
	fi
	go clean

.PHONY: all test deps devel-deps lint cover build clean
