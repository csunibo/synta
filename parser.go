package synta

import (
	"errors"
	"fmt"
	"strings"
)

func ParseSynta(contents string) (s Synta, err error) {
	lines := strings.Split(contents, "\n")

	var (
		i               = 0
		id              = Identifier("")
		def             = Definition{}
		definitionLines = []string{}
		filenameLine    = ""
	)
	if len(lines) > 1 {
		definitionLines = lines[:len(lines)-2]
		filenameLine = lines[len(lines)-1]
	} else {
		filenameLine = lines[0]
	}

	s.Definitions = map[Identifier]Definition{}
	for err == nil || i < len(definitionLines) {
		i, id, def, err = ParseNextDefinition(definitionLines, i)
		if _, ok := s.Definitions[id]; ok {
			return s, fmt.Errorf("Defintion for `%s` is provided twice", id)
		}
		s.Definitions[id] = def
	}
	if err != nil {
		return
	}

	s.Filename, err = ParseFilename(filenameLine)
	if len(s.Filename) == 0 {
		return s, errors.New("No filename construction specified")
	}

	return
}

func ParseNextDefinition(lines []string, start int) (i int, id Identifier, def Definition, err error) {
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
