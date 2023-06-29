package synta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkDefinitions(t *testing.T, definitions map[Identifier]Definition, expected map[string]string) {
	for id, def := range definitions {
		regexp, ok := expected[string(id)]
		if !ok {
			t.Errorf("Found unexpected defintion: %s = %s", string(id), def.Regexp.String())
		}

		assert.Equal(t, def.Regexp.String(), regexp)
		delete(expected, string(id))
	}

	if len(expected) > 0 {
		t.Errorf("Missing expected definitons: %v", expected)
	}
}

func TestParseSyntaWithEmptyFile(t *testing.T) {
	synta, err := ParseSynta("")
	assert.NotNil(t, err)
	assert.Empty(t, synta.Definitions)
	assert.Empty(t, synta.Filename)
}

func TestParseSyntaWithSingleDefinition(t *testing.T) {
	input := `test = a|b
> test-test.test`
	synta, err := ParseSynta(input)
	assert.Nil(t, err)
	assert.NotEmpty(t, synta.Definitions)

	exp := map[string]string{
		"test": "a|b",
	}
	checkDefinitions(t, synta.Definitions, exp)

	assert.NotEmpty(t, synta.Filename)
	assert.Equal(t, synta.Filename, Filename{
		Segments:  []Segment{{Identifier("test"), false}, {Identifier("test"), false}},
		Extension: Identifier("test"),
	})
}
