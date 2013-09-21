export GOPATH=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

all: reqs bindata build clean

reqs:
	@echo "\x1b[31;1mGetting dependencies...\x1b[0m"
	go get github.com/rafikk/imagick/imagick


build:
	@echo "\x1b[31;1mBuilding...\x1b[0m"
	go build

runtests:
	@echo "\x1b[31;1mTesting...\x1b[0m"
	go test

test: reqs runtests