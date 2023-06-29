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

	s.Filename, err = ParseFilename(filenameLine)
	return
}

func ParseNextDefinition(lines []string, start int) (i int, id Identifier, def Definition, err error) {
	i = start

	// TODO:
	// loop until definition row, accumulate comments,
	// when a def is reached, simply parse it and check the regexp. Then return
	for _, line := range lines {
		i++
		if line[0] == ';' {
			def.Comments = append(def.Comments, line)
		} else {
			parsed_line := strings.Split(line, " = ")
			id = Identifier(parsed_line[0])
			err = nil
			def.Regexp, err = regexp.Compile(parsed_line[1])
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

func ParseFilename(line string) (def []Segment, err error) {
	// TODO:
	// 1. Check if the string starts with "> " and then trim it
	// 2. Implement the DFA and compute the segments
	if len(line) < 2 || line[:2] != "> " {
		err = errors.New("Not a Filename")
		return def, errors.New("Not a Filename")
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
	def = push(def, &seg)

	return
}
