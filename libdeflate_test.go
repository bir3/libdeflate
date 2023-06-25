package libdeflate

import (
	"encoding/binary"
	"testing"
)

func BGZF_bsize(data []byte) int {
	if len(data) >= 28 {
		return int(binary.LittleEndian.Uint16(data[16:18])) + 1
	}
	panic("too short bgzf block")
}

func TestBGZFDecompress(t *testing.T) {
	data := []byte{
		// block 1: text: "lib"
		0x1f, 0x8b, 0x08, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x06, 0x00, 0x42, 0x43,
		0x02, 0x00, 0x1e, 0x00, 0xcb, 0xc9, 0x4c, 0x02, 0x00, 0xcc, 0x3b, 0x0f, 0xa9, 0x03,
		0x00, 0x00, 0x00,
		// block 2: embedded bgzf EOF marker empty block
		0x1f, 0x8b, 0x08, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x06, 0x00, 0x42, 0x43,
		0x02, 0x00, 0x1b, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// block 3: text: "deflate"
		0x1f, 0x8b, 0x08, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x06, 0x00, 0x42, 0x43,
		0x02, 0x00, 0x22, 0x00, 0x4b, 0x49, 0x4d, 0xcb, 0x49, 0x2c, 0x49, 0x05, 0x00, 0x05,
		0x6b, 0xb0, 0xe2, 0x07, 0x00, 0x00, 0x00,
	}
	out := make([]byte, 1000)
	dc := new(Decompressor)
	err := dc.Init()
	//var err error
	if err != nil {
		t.Fatal(err)
	}
	offset := 0
	outstr := ""
	for {
		var n int
		if len(data) == offset {
			break
		}

		n, err = dc.GzipDecompress(out, data[offset:])
		if err != nil {
			t.Fatal(err)
		}
		outstr += string(out[0:n])

		offset += BGZF_bsize(data[offset:])
	}
	dc.Cleanup()
	expect := "libdeflate"
	if outstr != expect {
		t.Fatalf("got %s but expected %s", outstr, expect)
	}
}

func decompress(t *testing.T, data []byte) string {
	out := make([]byte, 1000)
	dc := new(Decompressor)
	err := dc.Init()
	//var err error
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("# len(data) = %d\n", len(data))
	n, err := dc.Decompress(out, data)
	if err != nil {
		t.Fatal(err)
	}
	return string(out[0:n])
}

func TestBGZFCompress(t *testing.T) {
	data := "libdeflate"
	cc := new(Compressor)
	level := 1
	cc.Init(level)
	out := make([]byte, 1000)
	n_out := cc.Compress(out, []byte(data))
	outstr := decompress(t, out[0:n_out])
	expect := data
	if outstr != expect {
		t.Fatalf("got %s but expected %s", outstr, expect)
	}
}
