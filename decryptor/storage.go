package decryptor

import "io"

// Extension is a type helper for file extensions
type Extension string

// Storage defines an interface for storing decrypted video clips
type Storage interface {
	Save(Clip, io.Reader, Extension) (string, error)
}
