package format

import (
	"sort"

	"github.com/csunibo/synta"
)

// Format formats a synta structure into a string which represnts the contents
// of the associated synta file. The definitions are sorted by name. and blank
// lines are let as per internal convention.
func Format(syntaFile synta.Synta) (code string) {
	definitions := []string{}
	for id := range syntaFile.Definitions {
		definitions = append(definitions, string(id))
	}
	sort.Strings(definitions)

	for _, id := range definitions {
		def := syntaFile.Definitions[synta.Identifier(id)]
		for _, comment := range def.Comments {
			code += "; " + comment + "\n"
		}
		code += string(id) + " = " + def.Regexp.String() + "\n\n"
	}

	code += "> "
	for i, segment := range syntaFile.Filename.Segments {
		if segment.Optional {
			code += "(-" + string(segment.Identifier) + ")?"
		} else {
			code += string(segment.Identifier)
		}

		if i != len(syntaFile.Filename.Segments)-1 && !syntaFile.Filename.Segments[i+1].Optional {
			code += "-"
		}
	}
	code += "." + string(syntaFile.Filename.Extension) + "\n"

	return
}
