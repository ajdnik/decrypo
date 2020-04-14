package pluralsight

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io"
)

const (
	str1v1 = "706c7572616c7369676874"
	str2v1 = "063f7a59a2b2859f4cbeee30d62eec1723a93ec5a35105a4b00138de5e8efa194c71df279d03df459e4d8027783a007eb901ff2034b3f503c3a7ca0e41cbbc90e89eee7e8b9ae21bb855443c7f4be72a1df6e637480b154172fd2a76f725c2febee43b70fc"
	str1v2 = "00bf7b553901ae60eb13d15b1bcf"
	str2v2 = "028d0799899a25844bb073fac13438e4637a409f2ced3ef6a0320bdf0a402aed0b7a8c04bd9300dc65cb861f08d69e204144d36726ecb6178dc0147bb5ecdf88d89ff2d5c48170aaaa74438a409c323ac5665c5cade89efd0267037cd8426692a0"
)

var (
	// ErrUnknownStream defines an error sent when the stream being decoded is not recognized
	ErrUnknownStream = errors.New("unknown encrypted stream")
)

// xorBuff decrypts a byte buffer using xor operations with two unique keys
func xorBuff(n, offset int, buf, key1, key2 []byte) {
	for i := 0; i < n; i++ {
		num := key1[(offset+i)%len(key1)] ^ key2[(offset+i)%len(key2)] ^ byte((offset+i)%251)
		buf[i] = buf[i] ^ num
	}
}

// videoDecryptorFactory generates a new decryption io.Reader implementation based on the stream
func videoDecryptorFactory(r io.Reader) (videoDecryptor, error) {
	buff := bufio.NewReader(r)
	d, err := buff.Peek(3)
	if err != nil {
		return nil, err
	}
	key1v1, _ := hex.DecodeString(str1v1)
	key2v1, _ := hex.DecodeString(str2v1)
	// decrypt the first 3 bytes
	xorBuff(3, 0, d, key1v1, key2v1)
	// valid mp4 header should start with 3 zero bytes
	if d[0] == 0 && d[1] == 0 && d[2] == 0 {
		return videoDecryptor{
			Reader: &buff,
			Buf1:   key1v1,
			Buf2:   key2v1,
			Offset: 0,
		}, nil
	}
	// reverse the previous decryption
	xorBuff(3, 0, d, key1v1, key2v1)
	key1v2, _ := hex.DecodeString(str1v2)
	key2v2, _ := hex.DecodeString(str2v2)
	// decrypt the first 3 bytes using second set of keys
	xorBuff(3, 0, d, key1v2, key2v2)
	if d[0] == 0 && d[1] == 0 && d[2] == 0 {
		return videoDecryptor{
			Reader: &buff,
			Buf1:   key1v2,
			Buf2:   key2v2,
			Offset: 0,
		}, nil
	}
	return nil, ErrUnknownStream
}

type videoDecryptor struct {
	Reader io.Reader
	Key1   []byte
	Key2   []byte
	Offset int
}

// Read implements an io.Reader interface used to decrypt Pluralsight's videos
func (d *videoDecryptor) Read(buf []byte) (int, error) {
	n, err := d.Reader.Read(buf)
	if err != nil {
		return n, err
	}
	xorBuff(n, d.Offset, buf, d.Key1, d.Key2)
	d.Offset += n
	return n, err
}
