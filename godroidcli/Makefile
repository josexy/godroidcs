.PHONY: all build_run build run fmt clean
.DEFAULT_GOAL := all

BINARY=godroidcli
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

fmt:
	@go fmt .
	@go vet .

build:
	@CGO_ENABLE=0 GOOS="${GOOS}" GOARCH="${GOARCH}" go build -ldflags="-s -w" -o "${BINARY}" .

build_run:	build
	@./${BINARY}

run:
	@go run main.go

clean:
	@if [ -f ${BINARY} ] ;then rm ${BINARY} ;fi

all: fmt build

