shell := bash

build:
	go build -o ./out/nexus-minimal ./

run:
	make build
	./out/nexus-minimal

registry:
	docker build . -t astma/nexus-minimal:latest
	docker push astma/nexus-minimal:latest
