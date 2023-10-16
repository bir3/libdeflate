
all:
	CGO_ENABLED=1 go test
	CGO_ENABLED=0 go test --tags=use_slow_gzip
