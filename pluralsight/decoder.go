package pluralsight

import (
	"io"

	"github.com/ajdnik/decrypo/decryptor"
)

// Decoder decrypts Pluralsight's course videos
type Decoder struct{}

// Decode builds a video decryption stream
func (d *Decoder) Decode(r io.Reader) io.Reader {
	dec := newVideoDecryptor(r)
	return &dec
}

// Extension returns the decoded file extension
func (d *Decoder) Extension() decryptor.Extension {
	return "mp4"
}
