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
	// remove blank lines
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
		if lines[i] == "" {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}

	var (
		consumed        = 0
		id              = Identifier("")
		def             = Definition{}
		definitionLines = []string{}
		filenameLine    = ""
	)
	if len(lines) > 1 {
		definitionLines = lines[:len(lines)-1]
		filenameLine = lines[len(lines)-1]
	} else if len(lines) == 1 {
		err = errors.New("Missing either the filename or defintions")
	} else {
		err = errors.New("Empty file provided")
		return
	}

	s.Definitions = map[Identifier]Definition{}
	for len(definitionLines) > 0 {
		consumed, id, def, err = parseFirstDefinition(definitionLines)
		definitionLines = definitionLines[consumed:]
		if err != nil {
			return
		}

		if _, ok := s.Definitions[id]; ok {
			err = fmt.Errorf("Defintion for `%s` is provided twice", id)
			return
		} else if err == nil {
			s.Definitions[id] = def
		}
	}

	s.Filename.Segments, s.Filename.Extension, err = parseFilename(filenameLine)
	if err != nil {
		return
	}
	requiredIdentifiers := []Identifier{s.Filename.Extension}
	for _, seg := range s.Filename.Segments {
		requiredIdentifiers = append(requiredIdentifiers, seg.Identifier)
	}
	for _, id := range requiredIdentifiers {
		if _, ok := s.Definitions[id]; !ok {
			err = fmt.Errorf("Missing definition for `%s`", id)
			return
		}
	}

	return
}

func MustSynta(contents string) Synta {
	s, err := ParseSynta(contents)
	if err != nil {
		panic(err)
	}

	return s
}

// ParseNextDefinition loops from the start line, until a definition is found.
// All the lines from start to the defintion must be comments. If the defintion
// identifier is not valid, we return an error, otherwise, the definition index,
// the definition identifier and the definition itself are returned.
func parseFirstDefinition(lines []string) (consumed int, id Identifier, def Definition, err error) {
	for _, line := range lines {
		consumed++
		if line[0] == ';' {
			def.Comments = append(def.Comments, strings.TrimSpace(line[1:]))
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

// parseFilename checks if the line starts with "> ", or errors otherwise.
// Then, it parses a list of segments from the line using a DFA. If an invalid
// char is found, an error is returned, otherwise the result is the list of
// prased defintions.
func parseFilename(line string) (def []Segment, ext Identifier, err error) {
	if len(line) < 2 || line[:2] != "> " {
		err = errors.New("Not a Filename")
		return
	}
	line = line[2:]
	state := State0
	seg := Segment{}
	clear(&seg)

	col := 0
	for col = 0; err == nil && col < len(line); col++ {
		c := line[col]
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
				seg.Optional = true
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
			} else if c == '(' {
				seg.Optional = true
				state = State2
			} else {
				err = errors.New("Expected either a - or a . or a (")
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

	// ensure that we stop on an accepting state
	if err == nil && state != State8 {
		err = fmt.Errorf("Stopped at a non-accepting state (was %d, expected 8)", state)
	}
	// handle the filename extension
	ext = seg.Identifier

	// add debug information to the error string
	if err != nil {
		err = fmt.Errorf("Invalid char at column %d:\n%s\n%s\n%v", col, line, strings.Repeat(" ", col-1)+"^", err)
	}
	return
}
