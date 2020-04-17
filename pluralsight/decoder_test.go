package pluralsight_test

import (
	"testing"

	"github.com/ajdnik/decrypo/pluralsight"
)

func TestDecoder_Extension(t *testing.T) {
	decoder := pluralsight.Decoder{}

	ext := decoder.Extension()

	if ext != "mp4" {
		t.Errorf("got %v, want %v", ext, "mp4")
	}
}
