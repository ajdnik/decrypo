package file_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ajdnik/decrypo/decryptor"
	"github.com/ajdnik/decrypo/file"
)

var encodeTests = []struct {
	desc string
	in   []decryptor.Caption
	out  string
}{
	{"empty captions slice", []decryptor.Caption{}, ""},
	{"nil captions slice", nil, ""},
	{"sample caption 1", []decryptor.Caption{
		{StartMs: 0, EndMs: 0, Text: "hello", Clip: nil},
	}, fmt.Sprintf("1%[1]s00:00:00,000 --> 00:00:00,000%[1]shello%[1]s%[1]s", file.NewLine)},
	{"sample caption 2", []decryptor.Caption{
		{StartMs: 0, EndMs: 100, Text: "hello 2", Clip: nil},
	}, fmt.Sprintf("1%[1]s00:00:00,000 --> 00:00:00,100%[1]shello 2%[1]s%[1]s", file.NewLine)},
	{"ordered captions", []decryptor.Caption{
		{StartMs: 1200, EndMs: 3200, Text: "second", Clip: nil},
		{StartMs: 200, EndMs: 700, Text: "first", Clip: nil},
	}, fmt.Sprintf("1%[1]s00:00:00,200 --> 00:00:00,700%[1]sfirst%[1]s%[1]s2%[1]s00:00:01,200 --> 00:00:03,200%[1]ssecond%[1]s%[1]s", file.NewLine)},
}

func TestSrtEncoder_Encode(t *testing.T) {
	encoder := file.SrtEncoder{}
	for _, tt := range encodeTests {
		t.Run(tt.desc, func(t *testing.T) {
			r := encoder.Encode(tt.in)
			b, err := ioutil.ReadAll(r)
			if err != nil {
				t.Errorf("got error while reading output %v", err)
			}
			if string(b) != tt.out {
				t.Errorf("got '%v', want '%v'", string(b), tt.out)
			}
		})
	}
}
