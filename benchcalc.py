#! /usr/bin/env python3

import sys
import subprocess

MB = 1024 * 1024

# expect line like:
# BenchmarkClinvarDecompress-8   	       4	 323714698 ns/op
#

res = subprocess.run("bgzip -c -d dev/clinvar.vcf.gz|wc -c", shell=True, check=True, text=True, capture_output=True)
nbytes = int(res.stdout.strip())
assert nbytes > 800 * MB


def perf(f):
    for line in open(f):
        line = line.strip()
        e = line.split()
        if line.startswith("BenchmarkClinvarDecompress"):
            assert e[3] == "ns/op"
            dt_sec = int(e[2]) / 1e9
            name = f[len("tmp.") :]
            print(f"{line} # ", "%4.3f sec" % dt_sec, " %6.1f MB/sec" % (nbytes / dt_sec / 1e6), name)


for f in sys.argv[1:]:
    perf(f)
