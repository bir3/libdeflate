
# Go bindings for libdeflate 

```
# embeds libdeflate version 1.18
go get github.com/bir3/libdeflate@0.1.118
```

# Acknowledgments

The code was taken from https://github.com/grailbio/base/compress/libdeflate
and the libdeflate library updated.

Modifications:
- A build constraint for arm64 was removed.  The code runs fine on apple m1 (arm64)
- Avoid accidentially using slow stdlib gzip/zlib/deflate.  Now needs explict build tag `disable_libdefldate`
- added `actualDecompressor.Multistream(false)` to align libdeflate and stdlib wrappers

- https://github.com/ebiggers/libdeflate
- https://github.com/grailbio/base/compress/libdeflate `@ 0d762ae / 2023-04-14`

