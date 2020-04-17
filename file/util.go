package file

import (
	"os"
	"regexp"
	"strings"
)

var (
	isAbsWinDrive = regexp.MustCompile("^[a-zA-Z]\\:\\\\")
)

// ToUNC converts windows path to UNC path
func ToUNC(path string) string {
	// UNC can not use /
	path = strings.Replace(path, "/", "\\", -1)
	// if prefix starts with \\ we already have UNC path or server
	if strings.HasPrefix(path, "\\\\") {
		if strings.HasPrefix(path, "\\\\?\\") {
			return path
		}
		return "\\\\?\\UNC\\" + strings.TrimPrefix(path, "\\\\")
	}

	if isAbsWinDrive.MatchString(path) {
		return "\\\\?\\" + path
	}

	return path
}

// Exists checks if file exists on filesystem
func Exists(name string) bool {
	info, err := os.Lstat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
