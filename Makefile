VERSION?=dev
COMMIT=$(shell git rev-parse HEAD | cut -c -8)

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)"
MODFLAGS=-mod=vendor
TESTFLAGS=-cover

PLATFORMS:=darwin linux

dev:
	@echo building dist/ensong
	@go build ${MODFLAGS} ${LDFLAGS} -o dist/ensong ./cmd/ensong

prod: $(PLATFORMS)

$(PLATFORMS):
	@echo building dist/ensong-$@-amd64
	@GOOS=$@ GOARCH=amd64 go build ${MODFLAGS} ${LDFLAGS} -o dist/ensong-$@-amd64 ./cmd/ensong

test:
	@go test ${MODFLAGS} ${TESTFLAGS} ./...

clean:
	@rm -fr dist

.DEFAULT_GOAL := dev
.PHONY: dev prod $(PLATFORMS) test clean
