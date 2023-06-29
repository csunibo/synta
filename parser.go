package synta

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ParseSynta attempts to parse a file's contents into a Synta internal
// representation. If an error is encountered the parsing is aborted and the
// error returned
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
		definitionLines = lines[:len(lines)-1]
		filenameLine = lines[len(lines)-1]
	} else if len(lines) == 1 {
		filenameLine = lines[0]
	} else {
		err = errors.New("Empty file provided")
		return
	}

	s.Definitions = map[Identifier]Definition{}
	for i < len(definitionLines) {
		i, id, def, err = ParseNextDefinition(definitionLines, i)
		if _, ok := s.Definitions[id]; ok {
			err = fmt.Errorf("Defintion for `%s` is provided twice", id)
			return
		} else if err == nil {
			s.Definitions[id] = def
		}
	}
	if err != nil {
		return
	}

	s.Filename.Segments, s.Filename.Extension, err = ParseFilename(filenameLine)
	return
}

// ParseNextDefinition loops from the start line, until a definition is found.
// All the lines from start to the defintion must be comments. If the defintion
// identifier is not valid, we return an error, otherwise, the definition index,
// the definition identifier and the definition itself are returned.
func ParseNextDefinition(lines []string, start int) (i int, id Identifier, def Definition, err error) {
	i = start

	for _, line := range lines {
		i++
		if line[0] == ';' {
			def.Comments = append(def.Comments, line)
		} else {
			parsed_line := strings.Split(line, " = ")
			raw_id, expr := parsed_line[0], parsed_line[1]
			if !IdentifierRegexp.Match([]byte(raw_id)) {
				err = fmt.Errorf("Invalid identifier: %s", raw_id)
				return
			}
			id = Identifier(raw_id)
			def.Regexp, err = regexp.Compile(expr)
			return
		}
	}
	err = errors.New("No next definition")
	return
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

func push(segments []Segment, seg *Segment) (updatedSegments []Segment) {
	updatedSegments = append(segments, *seg)
	clear(seg)
	return
}

// ParseFilename checks if the line starts with "> ", or errors otherwise.
// Then, it parses a list of segments from the line using a DFA. If an invalid
// char is found, an error is returned, otherwise the result is the list of
// prased defintions.
func ParseFilename(line string) (def []Segment, ext Identifier, err error) {
	if len(line) < 2 || line[:2] != "> " {
		err = errors.New("Not a Filename")
		return
	}
	line = line[2:]
	state := State0
	seg := Segment{}
	clear(&seg)

	for i := 0; err == nil && i < len(line); i++ {
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
				def = push(def, &seg)
				state = State0
			} else if c == '(' {
				def = push(def, &seg)
				state = State2
			} else if c == '.' {
				def = push(def, &seg)
				state = State7
			} else {
				err = errors.New("Expected either a char, or a -, or a ( or a .")
			}
		case State2:
			if c == '-' {
				state = State3
			} else {
				err = errors.New("Expected a -")
			}
		case State3:
			if isLetter(c) {
				concat(&seg, c)
				state = State4
			} else {
				err = errors.New("Expected a char")
			}
		case State4:
			if isLetter(c) {
				concat(&seg, c)
			} else if c == ')' {
				state = State5
			} else {
				err = errors.New("Expected a char")
			}
		case State5:
			if c == '?' {
				def = push(def, &seg)
				state = State6
			} else {
				err = errors.New("Expected a ?")
			}
		case State6:
			if c == '-' {
				state = State0
			} else if c == '.' {
				state = State7
			} else {
				err = errors.New("Expected a - or a ?")
			}
		case State7:
			if isLetter(c) {
				concat(&seg, c)
				state = State8
			} else {
				err = errors.New("Expected a char")
			}
		case State8:
			if isLetter(c) {
				concat(&seg, c)
			} else {
				err = errors.New("Expected a char")
			}
		}
	}
	// handle the filename extension
	ext = seg.Identifier

	return
}
