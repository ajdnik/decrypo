package pluralsight

import (
	"database/sql"
	"encoding/json"
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

// captionEntry structure is a JSON object which is stored in the database and represents a single caption for a clip
type captionEntry struct {
	RelativeStartTime int    `json:"relativeStartTime"`
	RelativeEndTime   int    `json:"relativeEndTime"`
	Text              string `json:"text"`
}

// getCaptionsForClip retrieves clip captions and parses them into a struct
func getCaptionsForClip(clipID int, clip *decryptor.Clip, db *sql.DB) error {
	raw, err := db.Query(fmt.Sprintf("select ZCAPTIONS from ZCLIPCAPTIONSCD where ZCLIP=%v and ZLANGUAGECODE='en'", clipID))
	if err != nil {
		return err
	}
	defer raw.Close()
	if !raw.Next() {
		return nil
	}
	var data []byte
	err = raw.Scan(&data)
	if err != nil {
		return err
	}
	var entries []captionEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return err
	}
	// sort captions by start time
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].RelativeStartTime < entries[j].RelativeStartTime
	})
	for _, entry := range entries {
		caption := decryptor.Caption{
			StartMs: entry.RelativeStartTime,
			EndMs:   entry.RelativeEndTime,
			Text:    entry.Text,
			Clip:    clip,
		}
		clip.Captions = append(clip.Captions, caption)
	}
	return nil
}

// getClipsForModule retrieves video clips from an sqlite database that belong to a module
func getClipsForModule(modID int, mod *decryptor.Module, db *sql.DB) error {
	raw, err := db.Query(fmt.Sprintf("select Z_PK, ZTITLE, ZID from ZCLIPCD where ZMODULE=%v order by Z_FOK_MODULE asc", modID))
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
func getModulesForCourse(cID int, c *decryptor.Course, db *sql.DB) error {
	raw, err := db.Query(fmt.Sprintf("select Z_PK, ZTITLE from ZMODULECD where ZCOURSE=%v order by Z_FOK_COURSE asc", cID))
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
	raw, err := db.Query("select Z_PK, ZTITLE from ZCOURSEHEADERCD")
	if err != nil {
		return nil, err
	}
	defer raw.Close()
	courses := make([]decryptor.Course, 0)
	for raw.Next() {
		var id int
		var title string
		err = raw.Scan(&id, &title)
		if err != nil {
			return courses, err
		}
		course := decryptor.Course{
			Title:   title,
			Modules: make([]decryptor.Module, 0),
		}
		err = getModulesForCourse(id, &course, db)
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
	raw, err := db.Query("select count(*) from ZCLIPCD")
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
