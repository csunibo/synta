package regexp

import (
	"fmt"
	"regexp"

	"github.com/csunibo/synta"
)

func Convert(synta synta.Synta) (expr *regexp.Regexp, err error) {
    var finalString string = ""
    for i, segment := range synta.Filename.Segments {
        definition, isPresent := synta.Definitions[segment.Identifier]

        if !isPresent {
            err = fmt.Errorf("Missing definition for %s", segment.Identifier)
            return
        }

        if segment.Optional {
            finalString += "(-" + definition.Regexp.String() + ")?" 
        } else {
            finalString += "(" + definition.Regexp.String() + ")"
        }

        if i != len(synta.Filename.Segments) - 1 {
            finalString += "-"
        }
    }

    finalString += "\\.(" + synta.Definitions[synta.Filename.Extension].Regexp.String() + ")"
    expr, err = regexp.Compile(finalString)

    // Simplify when we use regexp/syntax
    // if err == nil {
    //     expr = expr.Simplify()
    // }
	return
}
