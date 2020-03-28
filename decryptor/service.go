package decryptor

// OnDecrypted defines a callback function called when a video clip is decrypted
type OnDecrypted func(Clip, *string)

// Service represents the decryption service which decrypts video courses
// and stores the courses in a readable format
type Service struct {
	Decoder        Decoder
	Storage        Storage
	CaptionEncoder CaptionEncoder
	Courses        CourseRepository
	Clips          ClipRepository
}

// DecryptAll decrypts all of the video courses contained in the courses repository
func (s *Service) DecryptAll(evt OnDecrypted) error {
	courses, err := s.Courses.FindAll()
	if err != nil {
		return err
	}
	for _, course := range courses {
		for _, module := range course.Modules {
			for _, clip := range module.Clips {
				if !s.Clips.Exists(&clip) {
					if evt != nil {
						evt(clip, nil)
					}
					continue
				}
				r, err := s.Clips.GetContent(&clip)
				if err != nil {
					return err
				}
				dec := s.Decoder.Decode(r)
				file, err := s.Storage.Save(clip, dec, s.Decoder.Extension())
				if err != nil {
					return err
				}
				r.Close()
				if len(clip.Captions) > 0 {
					enc := s.CaptionEncoder.Encode(clip.Captions)
					_, err = s.Storage.Save(clip, enc, s.CaptionEncoder.Extension())
					if err != nil {
						return err
					}
				}
				if evt != nil {
					evt(clip, &file)
				}
			}
		}
	}
	return nil
}
