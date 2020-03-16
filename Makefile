build:
	cd api && make 
	cd processor && make

test:
	cd processor && make test-unit
	cd processor && make test-integration

gen-proto:
	cd api && make 

