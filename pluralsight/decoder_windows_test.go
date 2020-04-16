package pluralsight_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ajdnik/decrypo/pluralsight"
)

var decoderTests = []struct {
	desc string
	in   []byte
	out  []byte
}{
	{"key v1", []byte{118, 82, 13, 52, 161, 175, 137, 129, 110, 235, 198, 107, 182, 86, 144, 120}, []byte{0, 0, 0, 28, 102, 116, 121, 112, 77, 52, 86, 32, 0, 0, 0, 1}},
	{"key v2", []byte{2, 51, 126, 235, 210, 234, 244, 147, 197, 218, 156, 152, 214, 246, 54, 84}, []byte{0, 0, 0, 36, 102, 116, 121, 112, 109, 112, 52, 50, 0, 0, 0, 0}},
}

func TestDecoder_DecodeBuffer(t *testing.T) {
	decoder := pluralsight.Decoder{}
	for _, tt := range decoderTests {
		t.Run(tt.desc, func(t *testing.T) {
			r := bytes.NewReader(tt.in)
			dec, err := decoder.Decode(r)

			if err != nil {
				t.Errorf("got an error while decoding, %v", err)
			}

			buff := bytes.NewBuffer([]byte{})
			_, err = buff.ReadFrom(dec)

			if err != nil {
				t.Errorf("got an error while reading decoded stream, %v", err)
			}

			out := buff.Bytes()

			if !bytes.Equal(tt.out, out) {
				t.Errorf("got %v, want %v", out, tt.out)
			}
		})
	}
}

var decoderErrorTests = []struct {
	desc string
	in   []byte
	out  error
}{
	{"no error", []byte{118, 82, 13, 52, 161, 175, 137, 129, 110, 235, 198, 107, 182, 86, 144, 120}, nil},
	{"unknown stream", []byte{1, 2, 3, 4, 5}, pluralsight.ErrUnknownStream},
	{"undefined buffer", nil, io.EOF},
}

func TestDecoder_DecodeError(t *testing.T) {
	decoder := pluralsight.Decoder{}
	for _, tt := range decoderErrorTests {
		t.Run(tt.desc, func(t *testing.T) {
			r := bytes.NewReader(tt.in)
			_, err := decoder.Decode(r)
			if err != tt.out {
				t.Errorf("got %v, want %v", err, tt.out)
			}
		})
	}
}
