package pluralsight

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajdnik/decrypo/decryptor"
	"github.com/ajdnik/decrypo/file"
)

// ClipRepository fetches encrypted video clips stored on the filesystem
type ClipRepository struct {
	Path string
}

// GetContent fetches an encrypted video clip stored on the filesystem based on clip's id
func (r *ClipRepository) GetContent(clip *decryptor.Clip) (io.ReadCloser, error) {
	repID := strings.ReplaceAll(clip.ID, "-", "")
	f, err := os.Open(filepath.Join(r.Path, fmt.Sprintf("%v.psv", repID)))
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Exists checks weather a video clip file exists
func (r *ClipRepository) Exists(clip *decryptor.Clip) bool {
	repID := strings.ReplaceAll(clip.ID, "-", "")
	cPath := filepath.Join(r.Path, fmt.Sprintf("%v.psv", repID))
	return file.Exists(cPath)
}
