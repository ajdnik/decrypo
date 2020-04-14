package pluralsight

import (
	"io"

	"github.com/ajdnik/decrypo/decryptor"
)

// Decoder decrypts Pluralsight's course videos
type Decoder struct{}

// Decode builds a video decryption stream
func (d *Decoder) Decode(r io.Reader) (io.Reader, error) {
	dec, err := videoDecryptorFactory(r)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

// Extension returns the decoded file extension
func (d *Decoder) Extension() decryptor.Extension {
	return "mp4"
}
