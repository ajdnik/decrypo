package pluralsight

import (
	"os"
	"path/filepath"

	"github.com/ajdnik/decrypo/file"
)

var (
	unknown = ""
)

// GetClipPath returns a default path where Pluralsight desktop app stores encrypted video clips
func GetClipPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return unknown, err
	}
	return file.ToUNC(filepath.Join(home, "AppData\\Local\\Pluralsight\\courses")), nil
}

// GetDbPath returns a default path where Pluralsight desktop app stores its sqlite database
func GetDbPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return unknown, err
	}
	return filepath.Join(home, "AppData\\Local\\Pluralsight\\pluralsight.db"), nil
}
