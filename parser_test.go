package synta

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DefinitionString struct {
	Comments     []string
	RegexpString string
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

	def := map[Identifier]DefinitionString{}
	def["test"] = DefinitionString{RegexpString: regexp.MustCompile("a|b").String()}
	defParsed := map[Identifier]DefinitionString{}
	for id, d := range synta.Definitions {
		fmt.Printf("Checking for key `%s`\n", id)
		defParsed[id] = DefinitionString{d.Comments, d.Regexp.String()}
		assert.Equal(t, def[id], defParsed[id])
	}

	assert.NotEmpty(t, synta.Filename)
	// TODO: controllare anche che filename contenga le strutture appropriate
	seg := Segment{Identifier: "test", Optional: false}
	assert.Contains(t, synta.Filename, seg)
}
