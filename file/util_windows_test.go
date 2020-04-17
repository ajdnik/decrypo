package file_test

import (
	"testing"

	"github.com/ajdnik/decrypo/file"
)

var existsTests = []struct {
	in  string
	out bool
}{
	{"C:\\unknown\\path\\file.txt", false},
	{"C:\\Windows\\System32\\calc.exe", true},
}

func TestExists(t *testing.T) {
	for _, tt := range existsTests {
		t.Run(tt.in, func(t *testing.T) {
			res := file.Exists(tt.in)
			if res != tt.out {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}
