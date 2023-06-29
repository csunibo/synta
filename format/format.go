package format

import "github.com/csunibo/synta"

func Format(synta synta.Synta) (code string) {
	for id, def := range synta.Definitions {
		for _, comment := range def.Comments {
			code += "; " + comment + "\n"
		}
		code += string(id) + " = " + def.Regexp.String() + "\n\n"
	}

	code += "> "
	for i, segment := range synta.Filename.Segments {
		if segment.Optional {
			code += "(-" + string(segment.Identifier) + ")?"
		} else {
			code += string(segment.Identifier)
		}

		if i != len(synta.Filename.Segments)-1 {
			code += "-"
		}
	}
	code += "." + string(synta.Filename.Extension) + "\n"

	return
}
