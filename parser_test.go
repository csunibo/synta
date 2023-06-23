package synta

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	def := Definition{Comments: []string{}, Regexp: regexp.MustCompile("a|b")}
	assert.Contains(t, synta.Definitions, def)

	assert.NotEmpty(t, synta.Filename)
	// TODO: controllare anche che filename contenga le strutture appropriate
}
