package synta

import (
	"errors"
	"fmt"
	"strings"
)

func ParseSynta(contents string) (s Synta, err error) {
	lines := strings.Split(contents, "\n")
	definitionLines, filenameLine := lines[:len(lines)-2], lines[len(lines)-1]

	var (
		i   = uint(0)
		id  = Identifier("")
		def = Definition{}
	)
	for err == nil {
		i, id, def, err = ParseNextDefinition(definitionLines, i)
		if _, ok := s.Definitions[id]; ok {
			return s, fmt.Errorf("Defintion for `%s` is provided twice", id)
		}
		s.Definitions[id] = def
	}

	s.Filename, err = ParseFilename(filenameLine)
	if len(s.Filename) == 0 {
		return s, errors.New("No filename construction specified")
	}

	return
}

func ParseNextDefinition(lines []string, start uint) (i uint, id Identifier, def Definition, err error) {
	i = start

	// TODO:
	// loop until definition row, accumulate comments,
	// when a def is reached, simply prase it and check the regexp. Then return
	return
}

func ParseFilename(line string) (def []Segment, err error) {
	// TODO:
	// 1. Check if the string starts with "> " and then trim it
	// 2. Implement the DFA and compute the segments
	return
}
