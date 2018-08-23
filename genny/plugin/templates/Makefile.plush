GO_BIN ?= go

deps:
	$(GO_BIN) get -v github.com/gobuffalo/packr/packr
	$(GO_BIN) get -v -t ./...

install: deps
	packr
	$(GO_BIN) install -v .
	packr clean

test:
	$(GO_BIN) test -tags ${TAGS} ./...

ci-test: deps
	$(GO_BIN) test -tags ${TAGS} -race ./...
