package libdeflate

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"os"
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

type Bench struct {
	bgzfData []byte
	out      []byte
	dc       *Decompressor
}

var g Bench

func init() {
	f := os.Getenv("DATA") //dev/clinvar.vcf")
	if f != "" {
		var err error
		g.bgzfData, err = os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "env DATA=%s but failed to read file: %s\n", f, err)
			os.Exit(7)
		}
		if len(g.bgzfData) < 70*1e6 {
			fmt.Fprintf(os.Stderr, "test file too small, must be > 70 MB")
			os.Exit(8)
		}
		g.out = make([]byte, 2000*1e6) // 2 GB
		g.dc = new(Decompressor)
		err = g.dc.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "dc.Init failed - %s", err)
			os.Exit(9)
		}

	}
}

func BenchmarkClinvarDecompress(b *testing.B) {
	if g.bgzfData == nil {
		b.Fatalf("env DATA not set - need bgzf file")
	}

	for i := 0; i < b.N; i++ {

		//b.Log("AAAA")
		//b.Logf("g.out = %d", len(g.out))
		data := g.bgzfData
		offset := 0
		n := 0
		k := 0
		for len(data) > offset {
			nx, err := g.dc.GzipDecompress(g.out[n:], data[offset:])
			if err != nil {
				b.Fatalf("decompress failed - %s", err)
			}
			offset += BGZF_bsize(data[offset:])
			n += nx
			k++
		}
		if n < 700*1e6 || n > 1500*1e6 {
			b.Fatalf("unexpected output size %d MB - %d bytes", n/1e6, n)
		}
		if n > len(g.out)/2 {
			b.Fatalf("output buffer too small")
		}

		b.SetBytes(int64(n))

		if os.Getenv("DATA_MD5") != "" {
			b.Logf("md5 verify: %x", md5.Sum(g.out[0:n]))
			if os.Getenv("DATA_MD5") != fmt.Sprintf("%x", md5.Sum(g.out[0:n])) {
				b.Fatalf("md5 verify failed")
			}
		}
		//b.ReportMetric(n float64, unit string)
	}
}
