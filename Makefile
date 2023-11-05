
all:
	@echo H4sIACcCSGUAA42OUQ6AMAhD/zkFf87E6IWW1INweGkZ+isko3l0gDuUHa3wwa7m4QN7vhWx1OVo2VXeIAY0jA2oLKk/ICtv5gafHHZ75JtangJpFGtv6I5Ri9krXCCKGZdlnuxP4EBM3lBXvEDM/H/YAwO2aM40AQAA|base64 -d |gunzip
	@echo "make test"
	@echo
	@echo "DATA=dev/clinvar.vcf.gz make b    # run benchmark; file at https://ftp.ncbi.nlm.nih.gov/pub/clinvar/vcf_GRCh38/clinvar.vcf.gz"
test:
	CGO_ENABLED=1 go test
	CGO_ENABLED=0 go test --tags=use_slow_gzip

# 323860052 ns/op

b bench:
	go test --bench . |tee tmp.libdeflate
	CGO_ENABLED=0 go test --tags=use_slow_gzip --bench .  |tee tmp.stdlib
	@python3 -c 'print("*"*80)'
	@echo "* REPORT:"
	./benchcalc.py tmp.*
	@rm tmp.libdeflate tmp.stdlib
