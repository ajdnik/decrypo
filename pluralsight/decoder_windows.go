package pluralsight

import (
	"encoding/hex"
	"io"
)

const (
	str1 = "706c7572616c7369676874"
	str2 = "063f7a59a2b2859f4cbeee30d62eec1723a93ec5a35105a4b00138de5e8efa194c71df279d03df459e4d8027783a007eb901ff2034b3f503c3a7ca0e41cbbc90e89eee7e8b9ae21bb855443c7f4be72a1df6e637480b154172fd2a76f725c2febee43b70fc"
)

// newVideoDecryptor generates a new decryption io.Reader implementation
func newVideoDecryptor(r io.Reader) videoDecryptor {
	buf1, _ := hex.DecodeString(str1)
	buf2, _ := hex.DecodeString(str2)
	return videoDecryptor{
		Reader: r,
		Buf1:   buf1,
		Buf2:   buf2,
		Offset: 0,
	}
}

type videoDecryptor struct {
	Reader io.Reader
	Buf1   []byte
	Buf2   []byte
	Offset int
}

// Read implements an io.Reader interface used to decrypt Pluralsight's videos
func (d *videoDecryptor) Read(buf []byte) (int, error) {
	n, err := d.Reader.Read(buf)
	if err != nil {
		return n, err
	}
	for i := 0; i < n; i++ {
		num := d.Buf1[(d.Offset+i)%len(d.Buf1)] ^ d.Buf2[(d.Offset+i)%len(d.Buf2)] ^ byte((d.Offset+i)%251)
		buf[i] = buf[i] ^ num
	}
	d.Offset += n
	return n, err
}
