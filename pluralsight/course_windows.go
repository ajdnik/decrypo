package pluralsight

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/ajdnik/decrypo/decryptor"
	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	unknownCount = -1
	// ErrNegativeCaptionTime is returned when caption entries read from database contain
	// negative values for either StartTime or EndTime
	ErrNegativeCaptionTime = errors.New("captions have a negative start or end time")
)

// CourseRepository fetches video course info from an sqlite database
type CourseRepository struct {
	Path string
}

// getCaptionsForClip retrieves clip captions and parses them into a struct
func getCaptionsForClip(clipID int32, clip *decryptor.Clip, db *sql.DB) error {
	raw, err := db.Query(fmt.Sprintf("select StartTime, EndTime, Text from ClipTranscript where ClipId=%v", clipID))
	if err != nil {
		return err
	}
	defer raw.Close()
	for raw.Next() {
		var startMs sql.NullInt32
		var endMs sql.NullInt32
		var text sql.NullString
		err = raw.Scan(&startMs, &endMs, &text)
		if err != nil {
			return err
		}
		// If any of the read values are NULL ignore the whole caption entry
		if !startMs.Valid || !endMs.Valid || !text.Valid {
			continue
		}
		// Return error of StartTime or EndTime are negative numbers
		// They represent miliseconds since the start of the clip so they should always be positive or 0
		if startMs.Int32 < 0 || endMs.Int32 < 0 {
			return ErrNegativeCaptionTime
		}
		caption := decryptor.Caption{
			StartMs: uint64(startMs.Int32),
			EndMs:   uint64(endMs.Int32),
			Text:    text.String,
			Clip:    clip,
		}
		clip.Captions = append(clip.Captions, caption)
	}
	// sort captions by start time
	sort.Slice(clip.Captions, func(i, j int) bool {
		return clip.Captions[i].StartMs < clip.Captions[j].StartMs
	})
	return nil
}

// getClipsForModule retrieves video clips from an sqlite database that belong to a module
func getClipsForModule(modID int32, mod *decryptor.Module, db *sql.DB) error {
	raw, err := db.Query(fmt.Sprintf("select Id, Title, Name from Clip where ModuleId=%v order by ClipIndex asc", modID))
	if err != nil {
		return err
	}
	defer raw.Close()
	ord := 1
	for raw.Next() {
		var id sql.NullInt32
		var title sql.NullString
		var uid sql.NullString
		err = raw.Scan(&id, &title, &uid)
		if err != nil {
			return err
		}
		// If any of the read values are NULL ignore the whole clip
		if !id.Valid || !title.Valid || !uid.Valid {
			continue
		}
		clip := decryptor.Clip{
			Order:    ord,
			Title:    title.String,
			ID:       uid.String,
			Module:   mod,
			Captions: make([]decryptor.Caption, 0),
		}
		err = getCaptionsForClip(id.Int32, &clip, db)
		if err != nil {
			return err
		}
		mod.Clips = append(mod.Clips, clip)
		ord++
	}
	return nil
}

// getModulesForCourse retrieves course modules from an sqlite database that belong to a video course
func getModulesForCourse(cName string, c *decryptor.Course, db *sql.DB) error {
	raw, err := db.Query(fmt.Sprintf("select Id, Title, Name, AuthorHandle from Module where CourseName='%v' order by ModuleIndex asc", cName))
	if err != nil {
		return err
	}
	defer raw.Close()
	ord := 1
	for raw.Next() {
		var id sql.NullInt32
		var title sql.NullString
		var uid sql.NullString
		var author sql.NullString
		err = raw.Scan(&id, &title, &uid, &author)
		if err != nil {
			return err
		}
		// If any of the read values are NULL skip the whole module
		if !id.Valid || !title.Valid || !uid.Valid || !author.Valid {
			continue
		}
		module := decryptor.Module{
			Order:  ord,
			Title:  title.String,
			ID:     uid.String,
			Author: author.String,
			Clips:  make([]decryptor.Clip, 0),
			Course: c,
		}
		err = getClipsForModule(id.Int32, &module, db)
		if err != nil {
			return err
		}
		c.Modules = append(c.Modules, module)
		ord++
	}
	return nil
}

// FindAll finds all of the video courses in the Pluralsight's sqlite database
func (r *CourseRepository) FindAll() ([]decryptor.Course, error) {
	db, err := sql.Open("sqlite3", r.Path)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	raw, err := db.Query("select Name, Title from Course")
	if err != nil {
		return nil, err
	}
	defer raw.Close()
	courses := make([]decryptor.Course, 0)
	for raw.Next() {
		var uid sql.NullString
		var title sql.NullString
		err = raw.Scan(&uid, &title)
		if err != nil {
			return courses, err
		}
		// If any of the read values are NULL ignore the whole course
		if !uid.Valid || !title.Valid {
			continue
		}
		course := decryptor.Course{
			Title:   title.String,
			ID:      uid.String,
			Modules: make([]decryptor.Module, 0),
		}
		err = getModulesForCourse(uid.String, &course, db)
		if err != nil {
			return courses, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

// ClipCount returns the number of all video clips in the Pluralsight's database
func (r *CourseRepository) ClipCount() (int, error) {
	db, err := sql.Open("sqlite3", r.Path)
	if err != nil {
		return unknownCount, err
	}
	defer db.Close()
	raw, err := db.Query("select count(*) from Clip")
	if err != nil {
		return unknownCount, err
	}
	defer raw.Close()
	if !raw.Next() {
		return unknownCount, sql.ErrNoRows
	}
	var count int
	err = raw.Scan(&count)
	if err != nil {
		return unknownCount, err
	}
	return count, nil
}
