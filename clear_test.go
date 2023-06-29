package synta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	input := `; a test comment
test = a|b
teest = a|b
needless = c|d
> test-teest.teest`
	synta, err := ParseSynta(input)
	assert.Nil(t, err)
	synta = Clear(synta)

	exp := StringDefintions{
		"test":  {"a|b", []string{"a test comment"}},
		"teest": {"a|b", []string(nil)},
	}
	checkDefinitions(t, synta.Definitions, exp)
}
