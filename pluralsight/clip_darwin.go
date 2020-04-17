package pluralsight

import (
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

// GetContent fetches an encrypted video clip stored on the filesystem based on clip's id
func (r *ClipRepository) GetContent(clip *decryptor.Clip) (io.ReadCloser, error) {
	if clip == nil {
		return nil, ErrClipUndefined
	}
	repID := strings.ReplaceAll(clip.ID, "-", "")
	f, err := r.FileOpen(filepath.Join(r.Path, fmt.Sprintf("%v.psv", repID)))
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
	repID := strings.ReplaceAll(clip.ID, "-", "")
	cPath := filepath.Join(r.Path, fmt.Sprintf("%v.psv", repID))
	return r.FileExists(cPath), nil
}
