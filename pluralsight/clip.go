package pluralsight

import (
	"errors"
	"os"
)

var (
	// ErrClipUndefined is returned when the clip argument is nil
	ErrClipUndefined = errors.New("clip argument is nil")
	// ErrModuleUndefined is returned when the clip's module property is nil
	ErrModuleUndefined = errors.New("module property is nil")
	// ErrCourseUndefined is returned when the module's course property is nil
	ErrCourseUndefined = errors.New("course property is nil")
)

// Open is a type declaration for the os.Open function
type Open func(string) (*os.File, error)

// Exists is a type declaration for the file.Exists function
type Exists func(string) bool
