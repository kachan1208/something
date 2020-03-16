.DEFAULT_GOAL := docker-build

build:
	CGO_ENABLED=0 go build ./ -o ./service 

docker-build:
	docker build -t processor .

test-unit:
	go test -v -test.race -tag=unit ./...

test-integration:
	go test -v -test.race -tag=integration ./tests/... 

.PHONY: build docker-build dep test-unit test-processor-integration