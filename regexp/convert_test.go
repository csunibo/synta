package regexp

import (
	"testing"

	"github.com/csunibo/synta"
	"github.com/stretchr/testify/assert"
)

func TestConvertBasic(t *testing.T) {
    basicContent := `test = a|b
    > test-test.test`
    basicSynta, err := synta.ParseSynta(basicContent)
    assert.Nil(t, err)

	expr, err := Convert(basicSynta)
	assert.Nil(t, err)
    assert.Equal(t, "(a|b)-(a|b)\\.(a|b)", expr.String())
}

