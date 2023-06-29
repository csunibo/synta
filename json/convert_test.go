package json

import (
	"testing"

	"github.com/csunibo/synta"
	"github.com/stretchr/testify/assert"
)

type StringDefintions = map[string]StringDefinition
type StringDefinition struct {
	Regexp   string
	Comments []string
}

func checkSynta(t *testing.T, syn synta.Synta, expected Synta) {
	assert.Equal(t, len(syn.Definitions), len(expected.Definitions))

	for id, def := range syn.Definitions {
		if len(def.Comments) == 0 {
			assert.Equal(t, len(def.Comments), len(expected.Definitions[string(id)].Comments))
		} else {
			assert.Equal(t, def.Comments, expected.Definitions[string(id)].Comments)
		}
		assert.Equal(t, def.Regexp.String(), expected.Definitions[string(id)].Regexp)
	}
	for i, e := range syn.Filename.Segments {
		assert.Equal(t, string(e.Identifier), string(expected.Filename.Segments[i].Identifier))
		assert.Equal(t, e.Optional, expected.Filename.Segments[i].Optional)
	}
	assert.Equal(t, string(syn.Filename.Extension), string(expected.Filename.Extension))
}

func TestConvertBasic(t *testing.T) {
	input := `test = a|b
> test-test.test`
	syn, err := synta.ParseSynta(input)
	assert.Nil(t, err)

	expectedConvert := Convert(syn)
	assert.Nil(t, err)
	checkSynta(t, syn, expectedConvert)
}
