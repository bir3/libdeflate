
# Go bindings for libdeflate 

```
# embeds libdeflate version 1.19
go get github.com/bir3/libdeflate@v0.4.119
```

# Example

```go
package main

import (
	"fmt"

	"github.com/bir3/libdeflate"
)

func main() {
	data := []byte{
		// text: "bgzf"
		0x1f, 0x8b, 0x08, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x06, 0x00, 0x42, 0x43,
		0x02, 0x00, 0x1f, 0x00, 0x4b, 0x4a, 0xaf, 0x4a, 0x03, 0x00, 0x20, 0x68, 0xf2, 0x8c,
		0x04, 0x00, 0x00, 0x00,
	}
	out := make([]byte, 1000)
	dc := new(libdeflate.Decompressor)
	err := dc.Init()
	if err != nil {
		panic(err)
	}

	n, err := dc.GzipDecompress(out, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(out[0:n]))
	dc.Cleanup()
}

```


# Acknowledgments

The code was taken from https://github.com/grailbio/base/tree/master/compress/libdeflate
and the [libdeflate](https://github.com/ebiggers/libdeflate) library updated.

Modifications:
- A build constraint for arm64 was removed.  The code runs fine on apple m1 (arm64)
- Avoid accidentially using slow stdlib gzip/zlib/deflate.  Now needs explict build tag `disable_libdefldate`
- added `actualDecompressor.Multistream(false)` to align libdeflate and stdlib wrappers

## code embedded/derived from:

- https://github.com/ebiggers/libdeflate
- https://github.com/grailbio/base/tree/0d762ae/compress/libdeflate  `@ 0d762ae / 2023-04-14`

