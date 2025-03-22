.DEFAULT_GOAL := build

init:
	git submodule update --init --recursive

docker:
	docker build -t "HPC-Scheduler" .

docker-run:
	docker run -d --restart always -p 127.0.0.1:8080:8080 "HPC-Scheduler"

coordinator:
	go build -v -o out/coodinator coodinator/cmd/main.go

build: coordinator tools cp-to-pwd

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
	go build -v -o out/tools/requester tools/requester/main.go
	go build -v -o out/tools/cli tools/cli/main.go
	cp out/tools/cli hpc

tool: tools

i: install

install:
	cp out/tools/cli /usr/bin/hpc

clean:
	rm -rf out

.PHONY: init docker docker-run build build-debug run tester clean tools tool