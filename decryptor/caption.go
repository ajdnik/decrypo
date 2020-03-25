package decryptor

import "io"

// Caption represents a video clip caption
type Caption struct {
	StartMs int
	EndMs   int
	Text    string
	Clip    *Clip
}

// CaptionEncoder defines an interface for encoding captions into different formats
type CaptionEncoder interface {
	Encode([]Caption) io.Reader
	Extension() Extension
}
