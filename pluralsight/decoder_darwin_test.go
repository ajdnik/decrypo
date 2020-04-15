package pluralsight_test

import (
	"bytes"
	"testing"

	"github.com/ajdnik/decrypo/pluralsight"
)

var decoderTests = []struct {
	desc string
	in   []byte
	out  []byte
}{
	{"sample 1", []byte{101, 101, 101, 65, 3, 17, 28, 21, 8, 21, 81, 87, 101, 101, 101, 101}, []byte{0, 0, 0, 36, 102, 116, 121, 112, 109, 112, 52, 50, 0, 0, 0, 0}},
}

func TestDecoder_Decode(t *testing.T) {
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
