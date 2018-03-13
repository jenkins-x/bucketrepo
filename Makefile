shell := bash

build:
	go build -o ./out/nexus-minimal ./

run:
	make build
	./out/nexus-minimal

registry:
	docker build . -t atsma/nexus-minimal
	docker push astma/nexus-minimal
