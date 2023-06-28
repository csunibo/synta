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

func ParseFilename(line string) (def []Segment, err error) {
	// TODO:
	// 1. Check if the string starts with "> " and then trim it
	// 2. Implement the DFA and compute the segments
	if line[:2] != "> " {
		return def, errors.New("Not a Filename")
	}

	line = line[2:]
	for i := 1; i < len(line); i++ {
		var seg = Segment{}
		switch { // Stato 0
		case line[i] >= 'a' && line[i] <= 'z':
			// concat
			// Passo in stato 1 e ciclo concat fino a '-'
			for ; line[i] >= 'a' && line[i] <= 'z'; i++ {
				// concat
			}
			if line[i] != '-' {
				// Errore
			} else {
				// list
				i++
			}
		case line[i] == '(':
			// Passo in stato 2 e gestisco il caso di parte ?
			i++
			if line[i] != '-' {
				// Error
			} else {
				i++
			}
			for ; line[i] >= 'a' && line[i] <= 'z'; i++ {
				// concat
			}
			if line[i] != ')' {
				// Error
			} else {
				i++
			}
			if line[i] != '?' {
				// Error
			} else {
				// list con opt == True
				i++
			}

		case line[i] == '.':
			// Passo in stato 6 per estensione file
			i++
			for ; line[i] >= 'a' && line[i] <= 'z'; i++ {
				// concat
			}
			if line[i] != '\n' {
				// Error
			}
		default:
			// Errore
		}
	}

	return def, err
}
