package pluralsight

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/ajdnik/decrypo/decryptor"
	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	unknownCount = -1
)

// CourseRepository fetches video course info from an sqlite database
type CourseRepository struct {
	Path string
}

// getCaptionsForClip retrieves clip captions and parses them into a struct
func getCaptionsForClip(clipID int, clip *decryptor.Clip, db *sql.DB) error {
	raw, err := db.Query(fmt.Sprintf("select StartTime, EndTime, Text from ClipTranscript where ClipId=%v", clipID))
	if err != nil {
		return err
	}
	defer raw.Close()
	for raw.Next() {
		var startMs int
		var endMs int
		var text string
		err = raw.Scan(&startMs, &endMs, &text)
		if err != nil {
			return err
		}
		caption := decryptor.Caption{
			StartMs: startMs,
			EndMs:   endMs,
			Text:    text,
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
func getClipsForModule(modID int, mod *decryptor.Module, db *sql.DB) error {
	raw, err := db.Query(fmt.Sprintf("select Id, Title, Name from Clip where ModuleId=%v order by ClipIndex asc", modID))
	if err != nil {
		return err
	}
	defer raw.Close()
	ord := 1
	for raw.Next() {
		var id int
		var title string
		var uid string
		err = raw.Scan(&id, &title, &uid)
		if err != nil {
			return err
		}
		clip := decryptor.Clip{
			Order:    ord,
			Title:    title,
			ID:       uid,
			Module:   mod,
			Captions: make([]decryptor.Caption, 0),
		}
		err = getCaptionsForClip(id, &clip, db)
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
	raw, err := db.Query(fmt.Sprintf("select Id, Title from Module where CourseName=%v order by ModuleIndex asc", cName))
	if err != nil {
		return err
	}
	defer raw.Close()
	ord := 1
	for raw.Next() {
		var id int
		var title string
		err = raw.Scan(&id, &title)
		if err != nil {
			return err
		}
		module := decryptor.Module{
			Order:  ord,
			Title:  title,
			Clips:  make([]decryptor.Clip, 0),
			Course: c,
		}
		err = getClipsForModule(id, &module, db)
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
		var name string
		var title string
		err = raw.Scan(&name, &title)
		if err != nil {
			return courses, err
		}
		course := decryptor.Course{
			Title:   title,
			Modules: make([]decryptor.Module, 0),
		}
		err = getModulesForCourse(name, &course, db)
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
