package format

import (
	"testing"

	"github.com/csunibo/synta"
	"github.com/stretchr/testify/assert"
)

func TestFormatBasic(t *testing.T) {
	basicContent := `def = a|b
test = c|d
> def.test
`
	basicSynta, err := synta.ParseSynta(basicContent)
	assert.Nil(t, err)

	formatted := Format(basicSynta)
	formattedContent := `def = a|b

test = c|d

> def.test
`
	assert.Equal(t, formattedContent, formatted)
}
