package decryptor_test

import (
	"io"
	"testing"

	"github.com/ajdnik/decrypo/decryptor"
)

type mockDecoder struct{}

func (md *mockDecoder) Decode(io.Reader) (io.Reader, error) {
	return nil, nil
}

func (md *mockDecoder) Extension() decryptor.Extension {
	return ""
}

type mockStorage struct{}

func (ms *mockStorage) Save(decryptor.Clip, io.Reader, decryptor.Extension) (string, error) {
	return "", nil
}

type mockCaptionEncoder struct{}

func (mce *mockCaptionEncoder) Encode([]decryptor.Caption) io.Reader {
	return nil
}

func (mce *mockCaptionEncoder) Extension() decryptor.Extension {
	return ""
}

type mockCourseRepository struct {
	NumClips int
}

func (mcr *mockCourseRepository) FindAll() ([]decryptor.Course, error) {
	course := decryptor.Course{
		Title: "",
		ID:    "",
		Modules: []decryptor.Module{
			{
				Order:  0,
				Title:  "",
				ID:     "",
				Author: "",
				Course: nil,
				Clips:  []decryptor.Clip{},
			},
		},
	}
	for i := 0; i < mcr.NumClips; i++ {
		c := decryptor.Clip{
			Order:    0,
			Title:    "",
			ID:       "",
			Module:   nil,
			Captions: []decryptor.Caption{},
		}
		course.Modules[0].Clips = append(course.Modules[0].Clips, c)
	}
	return []decryptor.Course{course}, nil
}

type mockReadCloser struct{}

func (mrc *mockReadCloser) Read([]byte) (int, error) {
	return 0, nil
}

func (mrc *mockReadCloser) Close() error {
	return nil
}

type mockClipRepository struct{}

func (mclr *mockClipRepository) GetContent(*decryptor.Clip) (io.ReadCloser, error) {
	return &mockReadCloser{}, nil
}

func (mclr *mockClipRepository) Exists(*decryptor.Clip) bool {
	return true
}

type mockCallback struct {
	Called int
}

func (mcb *mockCallback) Call(decryptor.Clip, *string) {
	mcb.Called++
}

func TestService_DecryptAllCallbackCalled(t *testing.T) {
	svc := decryptor.Service{
		Decoder:        &mockDecoder{},
		Storage:        &mockStorage{},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{6},
		Clips:          &mockClipRepository{},
	}

	cnt := mockCallback{}

	err := svc.DecryptAll(cnt.Call)

	if err != nil {
		t.Errorf("got an error %v", err)
	}

	if cnt.Called != 6 {
		t.Errorf("got %v, want %v", cnt.Called, 6)
	}
}
