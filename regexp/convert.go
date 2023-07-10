package regexp

import (
	"fmt"
	"regexp"

	"github.com/csunibo/synta"
)

func Convert(synta synta.Synta) (expr *regexp.Regexp, err error) {
	finalString, err := convertWithoutExtensionString(synta)
	if err != nil {
		return
	}

	finalString += "\\.(" + synta.Definitions[synta.Filename.Extension].Regexp.String() + ")"
    expr, err = regexp.Compile("^" + finalString + "$")

	// Simplify when we use regexp/syntax
	// if err == nil {
	//     expr = expr.Simplify()
	// }
	return
}

func convertWithoutExtensionString(synta synta.Synta) (expr string, err error) {
	for i, segment := range synta.Filename.Segments {
		definition, isPresent := synta.Definitions[segment.Identifier]

		if !isPresent {
			err = fmt.Errorf("Missing definition for %s", segment.Identifier)
			return
		}

		if segment.Optional {
			expr += "(-(" + definition.Regexp.String() + "))?"
		} else {
			expr += "(" + definition.Regexp.String() + ")"
		}

		if i != len(synta.Filename.Segments)-1 && !synta.Filename.Segments[i+1].Optional {
			expr += "-"
		}
	}
	return
}

func ConvertWithoutExtension(synta synta.Synta) (expr *regexp.Regexp, err error) {
    exp, err := convertWithoutExtensionString(synta)
    if err != nil {
        return
    }
    expr, err = regexp.Compile("^" + exp + "$")
	return
}
