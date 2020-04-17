package pluralsight_test

import (
	"os"
	"testing"

	"github.com/ajdnik/decrypo/decryptor"
	"github.com/ajdnik/decrypo/pluralsight"
)

type mockOpen struct {
	Arg string
}

func (mo *mockOpen) Open(f string) (*os.File, error) {
	mo.Arg = f
	return nil, nil
}

type mockExists struct {
	Arg string
}

func (me *mockExists) Exists(f string) bool {
	me.Arg = f
	return true
}

var getContentTests = []struct {
	desc string
	in   *decryptor.Clip
	arg  string
}{
	{"example 1", &decryptor.Clip{
		ID: "d6afc56e-daa3-4c03-91e3-8c9c9c915544",
	}, "/tmp/d6afc56edaa34c0391e38c9c9c915544.psv"},
	{"empty ID", &decryptor.Clip{}, "/tmp/.psv"},
}

func TestClipRepository_GetContent(t *testing.T) {
	open := mockOpen{}
	exists := mockExists{}
	repo := pluralsight.ClipRepository{
		Path:       "/tmp/",
		FileOpen:   open.Open,
		FileExists: exists.Exists,
	}
	for _, tt := range getContentTests {
		t.Run(tt.desc, func(t *testing.T) {
			repo.GetContent(tt.in)
			if open.Arg != tt.arg {
				t.Errorf("got %v, want %v", open.Arg, tt.arg)
			}
		})
	}
}

var getContentErrorsTest = []struct {
	desc string
	in   *decryptor.Clip
	err  error
}{
	{"clip exists", &decryptor.Clip{}, nil},
	{"no clip", nil, pluralsight.ErrClipUndefined},
}

func TestClipRepository_GetContentErrors(t *testing.T) {
	open := mockOpen{}
	exists := mockExists{}
	repo := pluralsight.ClipRepository{
		Path:       "/tmp/",
		FileOpen:   open.Open,
		FileExists: exists.Exists,
	}
	for _, tt := range getContentErrorsTest {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := repo.GetContent(tt.in)
			if err != tt.err {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}

var existsTests = []struct {
	desc string
	in   *decryptor.Clip
	arg  string
}{
	{"example 1", &decryptor.Clip{
		ID: "d6afc56e-daa3-4c03-91e3-8c9c9c915544",
	}, "/tmp/d6afc56edaa34c0391e38c9c9c915544.psv"},
	{"empty ID", &decryptor.Clip{}, "/tmp/.psv"},
}

func TestClipRepository_Exists(t *testing.T) {
	open := mockOpen{}
	exists := mockExists{}
	repo := pluralsight.ClipRepository{
		Path:       "/tmp/",
		FileOpen:   open.Open,
		FileExists: exists.Exists,
	}
	for _, tt := range getContentTests {
		t.Run(tt.desc, func(t *testing.T) {
			repo.Exists(tt.in)
			if exists.Arg != tt.arg {
				t.Errorf("got %v, want %v", exists.Arg, tt.arg)
			}
		})
	}
}

var existsErrorsTest = []struct {
	desc string
	in   *decryptor.Clip
	err  error
}{
	{"clip exists", &decryptor.Clip{}, nil},
	{"no clip", nil, pluralsight.ErrClipUndefined},
}

func TestClipRepository_ExistsErrors(t *testing.T) {
	open := mockOpen{}
	exists := mockExists{}
	repo := pluralsight.ClipRepository{
		Path:       "/tmp/",
		FileOpen:   open.Open,
		FileExists: exists.Exists,
	}
	for _, tt := range getContentErrorsTest {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := repo.Exists(tt.in)
			if err != tt.err {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}
