shell := bash

build:
	go build -o ./out/nexus-minimal ./

run:
	make build
	./out/nexus-minimal
