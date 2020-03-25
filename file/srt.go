package file

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/ajdnik/decrypo/decryptor"
)

// SrtEncoder builds an srt file from provided clip captions
type SrtEncoder struct{}

// msToString converts milliseconds to an srt time format
func msToString(ms int64) string {
	// convert milliseconds to nanoseconds
	var d time.Duration = time.Duration(ms * 1000000)
	hours := int64(d.Hours())
	minutes := int64(d.Minutes()) - hours*60
	seconds := int64(d.Seconds()) - int64(d.Minutes())*60
	milliseconds := d.Milliseconds() - int64(d.Seconds())*1000
	return fmt.Sprintf("%02d:%02d:%02d,%d", hours, minutes, seconds, milliseconds)
}

// Encode converts clip captions into an srt file
func (s *SrtEncoder) Encode(captions []decryptor.Caption) io.Reader {
	var sb strings.Builder
	// sort captions by start time
	sort.Slice(captions, func(i, j int) bool {
		return captions[i].StartMs < captions[j].StartMs
	})
	for idx, caption := range captions {
		sb.WriteString(fmt.Sprintf("%v%v%v --> %v%v%v%v%v", idx+1, NewLine, msToString(int64(caption.StartMs)), msToString(int64(caption.EndMs)), NewLine, caption.Text, NewLine, NewLine))
	}
	return strings.NewReader(sb.String())
}

// Extension returns the srt's file extension
func (s *SrtEncoder) Extension() decryptor.Extension {
	return "srt"
}
