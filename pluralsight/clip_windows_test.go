package pluralsight_test

import (
	"errors"
	"os"
	"testing"

	"github.com/ajdnik/decrypo/decryptor"
	"github.com/ajdnik/decrypo/pluralsight"
)

var (
	errMockOpen = errors.New("mock open error")
)

type mockOpen struct {
	Arg        string
	ThrowError bool
}

func (mo *mockOpen) Open(f string) (*os.File, error) {
	mo.Arg = f
	if mo.ThrowError {
		return nil, errMockOpen
	}
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
		Module: &decryptor.Module{
			ID:     "55197fb0-e440-473e-920a-cc785af5d82e",
			Author: "author",
			Course: &decryptor.Course{
				ID: "e0ef9f4a-60f7-46be-b3c4-4daa641c27c9",
			},
		},
	}, "C:\\Tmp\\e0ef9f4a-60f7-46be-b3c4-4daa641c27c9\\tSQ9XwTW4BXEQH+7SxZRmw==\\d6afc56e-daa3-4c03-91e3-8c9c9c915544.psv"},
	{"empty values", &decryptor.Clip{
		Module: &decryptor.Module{
			Course: &decryptor.Course{},
		},
	}, "C:\\Tmp\\uZg0vBm7rSRYCzrfoE+5Rw==\\.psv"},
}

func TestClipRepository_GetContent(t *testing.T) {
	open := mockOpen{}
	exists := mockExists{}
	repo := pluralsight.ClipRepository{
		Path:       "C:\\Tmp\\",
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
	desc     string
	in       *decryptor.Clip
	throwErr bool
	err      error
}{
	{"valid clip", &decryptor.Clip{
		Module: &decryptor.Module{
			Course: &decryptor.Course{},
		},
	}, false, nil},
	{"no clip", nil, false, pluralsight.ErrClipUndefined},
	{"clip without module", &decryptor.Clip{}, false, pluralsight.ErrModuleUndefined},
	{"module without course", &decryptor.Clip{
		Module: &decryptor.Module{},
	}, false, pluralsight.ErrCourseUndefined},
	{"open failed", &decryptor.Clip{
		Module: &decryptor.Module{
			Course: &decryptor.Course{},
		},
	}, true, errMockOpen},
}

func TestClipRepository_GetContentErrors(t *testing.T) {
	exists := mockExists{}
	repo := pluralsight.ClipRepository{
		Path:       "C:\\Tmp\\",
		FileExists: exists.Exists,
	}
	for _, tt := range getContentErrorsTest {
		t.Run(tt.desc, func(t *testing.T) {
			open := mockOpen{
				ThrowError: tt.throwErr,
			}
			repo.FileOpen = open.Open
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
		Module: &decryptor.Module{
			ID:     "55197fb0-e440-473e-920a-cc785af5d82e",
			Author: "author",
			Course: &decryptor.Course{
				ID: "e0ef9f4a-60f7-46be-b3c4-4daa641c27c9",
			},
		},
	}, "C:\\Tmp\\e0ef9f4a-60f7-46be-b3c4-4daa641c27c9\\tSQ9XwTW4BXEQH+7SxZRmw==\\d6afc56e-daa3-4c03-91e3-8c9c9c915544.psv"},
	{"empty values", &decryptor.Clip{
		Module: &decryptor.Module{
			Course: &decryptor.Course{},
		},
	}, "C:\\Tmp\\uZg0vBm7rSRYCzrfoE+5Rw==\\.psv"},
}

func TestClipRepository_Exists(t *testing.T) {
	open := mockOpen{}
	exists := mockExists{}
	repo := pluralsight.ClipRepository{
		Path:       "C:\\Tmp\\",
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
	{"valid clip", &decryptor.Clip{
		Module: &decryptor.Module{
			Course: &decryptor.Course{},
		},
	}, nil},
	{"no clip", nil, pluralsight.ErrClipUndefined},
	{"clip without module", &decryptor.Clip{}, pluralsight.ErrModuleUndefined},
	{"module without course", &decryptor.Clip{
		Module: &decryptor.Module{},
	}, pluralsight.ErrCourseUndefined},
}

func TestClipRepository_ExistsErrors(t *testing.T) {
	open := mockOpen{}
	exists := mockExists{}
	repo := pluralsight.ClipRepository{
		Path:       "C:\\Tmp\\",
		FileOpen:   open.Open,
		FileExists: exists.Exists,
	}
	for _, tt := range existsErrorsTest {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := repo.Exists(tt.in)
			if err != tt.err {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}
