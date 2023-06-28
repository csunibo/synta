package synta

import (
	"errors"
	"fmt"
	"regexp"
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
	// when a def is reached, simply parse it and check the regexp. Then return
	for i, line := range lines {
		if line[0] == ';' {
			def.Comments = append(def.Comments, line)
		} else {
			parsed_line := strings.Split(line, " = ")
			id := Identifier(parsed_line[0])
			def.Regexp = regexp.MustCompile(parsed_line[1])
			return i, id, def, nil
		}
	}
	return i, id, def, errors.New("No next Definition")
}

type State uint8

const (
	State0 State = iota
	State1
	State2
	State3
	State4
	State5
	State6
	State7
	State8
	State9
)

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func concat(seg *Segment, c byte) {
	seg.Identifier = Identifier(string(seg.Identifier) + string(c))
}

func clear(seg *Segment) {
	seg.Identifier = ""
	seg.Optional = false
}

func push(segments []Segment, seg *Segment) {
	segments = append(segments, *seg)
	clear(seg)
}

func ParseFilename(line string) (def []Segment, err error) {
	// TODO:
	// 1. Check if the string starts with "> " and then trim it
	// 2. Implement the DFA and compute the segments
	if line[:2] != "> " {
		return def, errors.New("Not a Filename")
	}
	line = line[2:]
	state := State0
	seg := Segment{}
	clear(&seg)

	for i := 1; err == nil && i < len(line); i++ {
		c := line[i]
		switch state {
		case State0:
			if isLetter(c) {
				concat(&seg, c)
				state = State1
			} else if c == '(' {
				state = State2
			} else {
				err = errors.New("Expected either a char or a (")
			}
		case State1:
			if isLetter(c) {
				concat(&seg, c)
			} else if c == '-' {
				push(def, &seg)
				state = State0
			} else {
				err = errors.New("Expected either a char or a -")
			}
		}
	}

	return def, err
}
