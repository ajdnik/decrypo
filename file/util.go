package file

import "os"

// Exists checks if file exists on filesystem
func Exists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
