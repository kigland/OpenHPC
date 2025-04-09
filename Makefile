.DEFAULT_GOAL := build

init:
	git submodule update --init --recursive

docker:
	docker build -t "OpenHPC" .

docker-run:
	docker run -d --restart always -p 127.0.0.1:8080:8080 "OpenHPC"

coordinator:
	go build -v -o out/coodinator coodinator/cmd/main.go

build: coordinator tools cp-to-pwd i

all: coordinator tools

cp-to-pwd:
	cp out/tools/cli hpc

publish:
	go build -v -gcflags="all=-N -l" -o out/serv cmd/serv/main.go

run:
	./build/serv

tester:
	mkdir -p out/tester
	go build -v -o out/tester/docker tester/docker/main.go
	go build -v -o out/tester/docker-jupyter tester/docker-jupyter/main.go

tools:
# go build -v -o out/tools/requester tools/requester/main.go
	go build -v -o out/tools/cli tools/cli/main.go

tool: tools

i: install

install:
	sudo cp out/tools/cli /usr/bin/hpc

clean:
	rm -rf out

api:
	bash gen_api.sh

.PHONY: init docker docker-run build build-debug run tester clean tools tool api