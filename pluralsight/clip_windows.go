package pluralsight

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/ajdnik/decrypo/decryptor"
)

// ClipRepository fetches encrypted video clips stored on the filesystem
type ClipRepository struct {
	Path       string
	FileOpen   Open
	FileExists Exists
}

// computeModuleHash generates a filesystem safe token from module data
func computeModuleHash(module *decryptor.Module) string {
	name := module.ID + "|" + module.Author
	h := md5.New()
	io.WriteString(h, name)
	enc := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return strings.ReplaceAll(enc, "/", "_")
}

// GetContent fetches an encrypted video clip stored on the filesystem
func (r *ClipRepository) GetContent(clip *decryptor.Clip) (io.ReadCloser, error) {
	if clip == nil {
		return nil, ErrClipUndefined
	}
	if clip.Module == nil {
		return nil, ErrModuleUndefined
	}
	if clip.Module.Course == nil {
		return nil, ErrCourseUndefined
	}
	f, err := r.FileOpen(filepath.Join(r.Path, clip.Module.Course.ID, computeModuleHash(clip.Module), fmt.Sprintf("%v.psv", clip.ID)))
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Exists checks weather a video clip file exists
func (r *ClipRepository) Exists(clip *decryptor.Clip) (bool, error) {
	if clip == nil {
		return false, ErrClipUndefined
	}
	if clip.Module == nil {
		return false, ErrModuleUndefined
	}
	if clip.Module.Course == nil {
		return false, ErrCourseUndefined
	}
	cPath := filepath.Join(r.Path, clip.Module.Course.ID, computeModuleHash(clip.Module), fmt.Sprintf("%v.psv", clip.ID))
	return r.FileExists(cPath), nil
}
