.DEFAULT_GOAL := docker-build

build:
	CGO_ENABLED=0 go build -o ./service -a ./ 

docker-build:
	docker build -t processor .

test-unit:
	go test -v -test.race -tags=unit -benchmem -bench=. ./... 

test-integration:
	go test -v -test.race --tags=integration ./tests/... 

.PHONY: bench build docker-build dep test-unit test-processor-integration
