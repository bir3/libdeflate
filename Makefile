
all:
	CGO_ENABLED=1 go test
	CGO_ENABLED=0 go test --tags=disable_libdeflate
