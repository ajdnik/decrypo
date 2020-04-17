package file_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/ajdnik/decrypo/decryptor"
	"github.com/ajdnik/decrypo/file"
)

func stubMkdirAll(string, os.FileMode) error {
	return nil
}

func stubWriteFile(string, []byte, os.FileMode) error {
	return nil
}

var storageTests = []struct {
	desc string
	clip decryptor.Clip
	buf  []byte
	ext  decryptor.Extension
	out  string
	err  error
}{
	{"correct path structure", decryptor.Clip{
		Order:    1,
		Title:    "clip",
		Captions: []decryptor.Caption{},
		Module: &decryptor.Module{
			Order: 1,
			Title: "module",
			Course: &decryptor.Course{
				Title: "course",
			},
		},
	}, []byte{1, 2, 3}, decryptor.Extension("ext"), "/path/course/1-module/1-clip.ext", nil},
	{"error if no course", decryptor.Clip{
		Order:    1,
		Title:    "clip",
		Captions: []decryptor.Caption{},
		Module: &decryptor.Module{
			Order:  1,
			Title:  "module",
			Course: nil,
		},
	}, []byte{1, 2, 3}, decryptor.Extension("ext"), "", file.ErrNil},
	{"error if no module", decryptor.Clip{
		Order:  1,
		Title:  "clip",
		Module: nil,
	}, []byte{1, 2, 3}, decryptor.Extension("ext"), "", file.ErrNil},
	{"correct path with long names", decryptor.Clip{
		Order: 1,
		Title: "clip with long name",
		Module: &decryptor.Module{
			Order: 1,
			Title: "module with long name",
			Course: &decryptor.Course{
				Title: "course with long name",
			},
		},
	}, []byte{1, 2, 3}, decryptor.Extension("ext"), "/path/course-with-long-name/1-module-with-long-name/1-clip-with-long-name.ext", nil},
}

func TestStorage_Save(t *testing.T) {
	storage := file.Storage{
		Path:      "/path/",
		MkdirAll:  stubMkdirAll,
		WriteFile: stubWriteFile,
	}
	for _, tt := range storageTests {
		t.Run(tt.desc, func(t *testing.T) {
			r := bytes.NewReader(tt.buf)
			fn, err := storage.Save(tt.clip, r, tt.ext)
			if fn != tt.out {
				t.Errorf("got %v, want %v", fn, tt.out)
			}
			if err != tt.err {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}
