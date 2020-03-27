package decryptor

import "io"

// Clip represents a video clip
type Clip struct {
	Order  int
	Title  string
	ID     string
	Module *Module
}

// ClipRepository defines an interface for fetching video clips
type ClipRepository interface {
	GetContent(*Clip) (io.ReadCloser, error)
	Exists(*Clip) bool
}
