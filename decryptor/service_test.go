package decryptor_test

import (
	"errors"
	"io"
	"testing"

	"github.com/ajdnik/decrypo/decryptor"
)

var (
	errMockCourses        = errors.New("mock courses error")
	errMockClips          = errors.New("mock clips error")
	errMockStorageDecoder = errors.New("mock storage error for decoders")
	errMockStorageEncoder = errors.New("mock storage error for encoders")
	errMockDecoder        = errors.New("mock decoder error")
	decoderExt            = decryptor.Extension("decoder")
	encoderExt            = decryptor.Extension("encoder")
)

type mockDecoder struct {
	ShowError bool
}

func (md *mockDecoder) Decode(io.Reader) (io.Reader, error) {
	if md.ShowError {
		return nil, errMockDecoder
	}
	return nil, nil
}

func (md *mockDecoder) Extension() decryptor.Extension {
	return decoderExt
}

type mockStorage struct {
	ShowDecoderError bool
	ShowEncoderError bool
}

func (ms *mockStorage) Save(c decryptor.Clip, r io.Reader, ext decryptor.Extension) (string, error) {
	if ms.ShowDecoderError && ext == decoderExt {
		return "", errMockStorageDecoder
	}
	if ms.ShowEncoderError && ext == encoderExt {
		return "", errMockStorageEncoder
	}
	return "", nil
}

type mockCaptionEncoder struct{}

func (mce *mockCaptionEncoder) Encode([]decryptor.Caption) io.Reader {
	return nil
}

func (mce *mockCaptionEncoder) Extension() decryptor.Extension {
	return encoderExt
}

type mockCourseRepository struct {
	NumClips  int
	ShowError bool
}

func (mcr *mockCourseRepository) FindAll() ([]decryptor.Course, error) {
	if mcr.ShowError {
		return nil, errMockCourses
	}
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
			Order:  0,
			Title:  "",
			ID:     "",
			Module: nil,
			Captions: []decryptor.Caption{
				{
					StartMs: 0,
					EndMs:   0,
					Text:    "",
					Clip:    nil,
				},
			},
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

type mockClipRepository struct {
	ShowError bool
	DoesExist bool
}

func (mclr *mockClipRepository) GetContent(*decryptor.Clip) (io.ReadCloser, error) {
	if mclr.ShowError {
		return nil, errMockClips
	}
	return &mockReadCloser{}, nil
}

func (mclr *mockClipRepository) Exists(*decryptor.Clip) bool {
	return mclr.DoesExist
}

type mockCallback struct {
	Called int
}

func (mcb *mockCallback) Call(decryptor.Clip, *string) {
	mcb.Called++
}

var decryptAllCallbacksTest = []struct {
	desc string
	in   decryptor.Service
	out  int
}{
	{"6 clips without errors", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 6, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, 6},
	{"7 clips do not exists", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 7, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: false},
	}, 7},
	{"clips error", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 1, ShowError: false},
		Clips:          &mockClipRepository{ShowError: true, DoesExist: true},
	}, 0},
	{"storage error for decoder", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: true},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 1, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, 0},
	{"storage error for encoder", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: true, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 1, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, 0},
	{"decoder error", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: true},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 1, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, 0},
}

func TestService_DecryptAllCallbacks(t *testing.T) {
	for _, tt := range decryptAllCallbacksTest {
		t.Run(tt.desc, func(t *testing.T) {
			cnt := mockCallback{}
			tt.in.DecryptAll(cnt.Call)
			if cnt.Called != tt.out {
				t.Errorf("got %v, want %v", cnt.Called, tt.out)
			}
		})
	}
}

var decryptAllErrorsTest = []struct {
	desc string
	in   decryptor.Service
	out  error
}{
	{"no errors", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 0, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, nil},
	{"no errors with 6 clips", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 6, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, nil},
	{"courses error", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 0, ShowError: true},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, errMockCourses},
	{"clips error not shown because no clips", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 0, ShowError: false},
		Clips:          &mockClipRepository{ShowError: true, DoesExist: true},
	}, nil},
	{"clips error", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 1, ShowError: false},
		Clips:          &mockClipRepository{ShowError: true, DoesExist: true},
	}, errMockClips},
	{"storage error for decoder", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: true},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 1, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, errMockStorageDecoder},
	{"storage error for encoder", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: false},
		Storage:        &mockStorage{ShowEncoderError: true, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 1, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, errMockStorageEncoder},
	{"decoder error", decryptor.Service{
		Decoder:        &mockDecoder{ShowError: true},
		Storage:        &mockStorage{ShowEncoderError: false, ShowDecoderError: false},
		CaptionEncoder: &mockCaptionEncoder{},
		Courses:        &mockCourseRepository{NumClips: 1, ShowError: false},
		Clips:          &mockClipRepository{ShowError: false, DoesExist: true},
	}, nil},
}

func TestService_DecryptAll(t *testing.T) {
	for _, tt := range decryptAllErrorsTest {
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.in.DecryptAll(nil)
			if err != tt.out {
				t.Errorf("got %v, want %v", err, tt.out)
			}
		})
	}
}
