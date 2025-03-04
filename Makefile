.DEFAULT_GOAL := build

init:
	git submodule update --init --recursive

docker:
	docker build -t "gin-template" .

docker-run:
	docker run -d --restart always -p 127.0.0.1:8080:8080 "gin-template"

build:
	bash build.sh

build-debug:
	go build -v -gcflags="all=-N -l" -o out/serv cmd/serv/main.go

run:
	./build/serv

.PHONY: init docker docker-run build build-debug run
