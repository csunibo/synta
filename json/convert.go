package json

import (
	"encoding/json"

	"github.com/csunibo/synta"
)

func Convert(synta synta.Synta) (s Synta) {
	s.Definitions = map[string]Definition{}
	for id, def := range synta.Definitions {
		if len(def.Comments) == 0 {
			s.Definitions[string(id)] = Definition{
				Comments: []string{}, Regexp: def.Regexp.String()}
		} else {
			s.Definitions[string(id)] = Definition{
				Comments: def.Comments, Regexp: def.Regexp.String()}
		}
	}
	for _, e := range synta.Filename.Segments {
		s.Filename.Segments = append(s.Filename.Segments, Segment{
			Identifier: string(e.Identifier), Optional: e.Optional})
	}
	s.Filename.Extension = string(synta.Filename.Extension)
	return
}

func ToJson(synta synta.Synta) (buf []byte, err error) {
	buf, err = json.Marshal(Convert(synta))
	return
}
