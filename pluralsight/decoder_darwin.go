package pluralsight

import "io"

// videoDecryptorFactory generates a new decryption io.Reader implementation
func videoDecryptorFactory(r io.Reader) (videoDecryptor, error) {
	return videoDecryptor{
		Reader: r,
	}, nil
}

type videoDecryptor struct {
	Reader io.Reader
}

// Read implements an io.Reader interface used to decrypt Pluralsight's videos
func (d *videoDecryptor) Read(buf []byte) (int, error) {
	n, err := d.Reader.Read(buf)
	if err != nil {
		return n, err
	}
	for i := 0; i < n; i++ {
		buf[i] = buf[i] ^ 101
	}
	return n, err
}
