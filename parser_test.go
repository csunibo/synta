package synta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StringDefintions = map[string]StringDefinition
type StringDefinition struct {
	Regexp   string
	Comments []string
}

func checkDefinitions(t *testing.T, definitions map[Identifier]Definition, expected map[string]StringDefinition) {
	for id, def := range definitions {
		regexp, ok := expected[string(id)]
		if !ok {
			t.Errorf("Found unexpected defintion: %s = %s", string(id), def.Regexp.String())
		}

		assert.Equal(t, regexp.Regexp, def.Regexp.String())
		assert.Equal(t, regexp.Comments, def.Comments)
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

func TestParseSyntaWithMissingFilename(t *testing.T) {
	_, err := ParseSynta("def = a|b")
	assert.NotNil(t, err)
}

func TestParseSyntaWithInvalidRegexp(t *testing.T) {
	_, err := ParseSynta("def = +")
	assert.NotNil(t, err)
}

func TestParseSyntaWithMissingDefinition(t *testing.T) {
	input := `def = a|b
> missingdef`
	_, err := ParseSynta(input)
	assert.NotNil(t, err)
}

func TestParseSyntaWithMissingExtensionDefinition(t *testing.T) {
	input := `def = a|b
> def.missingdef`
	_, err := ParseSynta(input)
	assert.NotNil(t, err)
}

// TODO: errore regexp, errore non esiste def, errore non esiste def nella ext

func TestParseSyntaWithSingleDefinition(t *testing.T) {
	input := `test = a|b
> test-test.test`
	synta, err := ParseSynta(input)
	assert.Nil(t, err)
	assert.NotEmpty(t, synta.Definitions)

	exp := StringDefintions{
		"test": {"a|b", []string(nil)},
	}
	checkDefinitions(t, synta.Definitions, exp)

	assert.NotEmpty(t, synta.Filename)
	id_test := Identifier("test")
	assert.Equal(t, synta.Filename, Filename{
		Segments: []Segment{
			{
				kind:        SegmentTypeIdentifier,
				value:       &id_test,
				subsegments: []Segment(nil),
			},
			{
				kind:        SegmentTypeIdentifier,
				value:       &id_test,
				subsegments: []Segment(nil),
			},
		},
		Extension: Identifier("test"),
	})
}

func TestParseSyntaWithSingleDefinitionSingleComment(t *testing.T) {
	input := `; a test comment
test = a|b
> test-test.test`
	synta, err := ParseSynta(input)
	assert.Nil(t, err)
	assert.NotEmpty(t, synta.Definitions)

	exp := StringDefintions{
		"test": {"a|b", []string{"a test comment"}},
	}
	checkDefinitions(t, synta.Definitions, exp)

	assert.NotEmpty(t, synta.Filename)
	id_test := Identifier("test")
	assert.Equal(t, synta.Filename, Filename{
		Segments:  []Segment{{SegmentTypeIdentifier, &id_test, []Segment(nil)}, {SegmentTypeIdentifier, &id_test, []Segment(nil)}},
		Extension: Identifier("test"),
	})
}

func TestParseSyntaWithSingleDefinitionMultipleComments(t *testing.T) {
	input := `; a test comment
; a second comment
test = a|b
> test-test.test`
	synta, err := ParseSynta(input)
	assert.Nil(t, err)
	assert.NotEmpty(t, synta.Definitions)

	exp := StringDefintions{
		"test": {"a|b", []string{"a test comment", "a second comment"}},
	}
	checkDefinitions(t, synta.Definitions, exp)

	assert.NotEmpty(t, synta.Filename)
	id_test := Identifier("test")
	assert.Equal(t, synta.Filename, Filename{
		Segments:  []Segment{{SegmentTypeIdentifier, &id_test, []Segment(nil)}, {SegmentTypeIdentifier, &id_test, []Segment(nil)}},
		Extension: Identifier("test"),
	})
}

func TestParseSyntaWithMultipleDefinitionMultipleComments(t *testing.T) {
	input := `; a test comment
; a second comment
test = a|b
; a test comment
; a second comment
teest = a|b
> test-teest.teest`
	synta, err := ParseSynta(input)
	assert.Nil(t, err)
	assert.NotEmpty(t, synta.Definitions)

	exp := StringDefintions{
		"test":  {"a|b", []string{"a test comment", "a second comment"}},
		"teest": {"a|b", []string{"a test comment", "a second comment"}},
	}
	checkDefinitions(t, synta.Definitions, exp)

	assert.NotEmpty(t, synta.Filename)
	id_test := Identifier("test")
	id_teest := Identifier("teest")
	assert.Equal(t, synta.Filename, Filename{
		Segments:  []Segment{{SegmentTypeIdentifier, &id_test, []Segment(nil)}, {SegmentTypeIdentifier, &id_teest, []Segment(nil)}},
		Extension: Identifier("teest"),
	})
}

func TestParseSyntaWithOptional(t *testing.T) {
	input := `; a test comment
; a second comment
test = a|b
> test(-test)?.test`
	synta, err := ParseSynta(input)
	assert.Nil(t, err)
	assert.NotEmpty(t, synta.Definitions)

	exp := StringDefintions{
		"test": {"a|b", []string{"a test comment", "a second comment"}},
	}
	checkDefinitions(t, synta.Definitions, exp)

	assert.NotEmpty(t, synta.Filename)
	id_test := Identifier("test")
	assert.Equal(t, synta.Filename, Filename{
		Segments: []Segment{
			{
				SegmentTypeIdentifier,
				&id_test,
				[]Segment(nil)},
			{
				SegmentTypeOptional,
				nil,
				[]Segment{{SegmentTypeIdentifier, &id_test, []Segment(nil)}},
			},
		},
		Extension: Identifier("test"),
	})

}

func TestParseSyntaWithNestedOptional(t *testing.T) {
	input := `; a test comment
; a second comment
test = a|b
> test(-test(-test)?-test(-test)?)?.test`
	synta, err := ParseSynta(input)
	assert.Nil(t, err)
	assert.NotEmpty(t, synta.Definitions)

	exp := StringDefintions{
		"test": {"a|b", []string{"a test comment", "a second comment"}},
	}
	checkDefinitions(t, synta.Definitions, exp)

	assert.NotEmpty(t, synta.Filename)

	id_test := Identifier("test")
	assert.Equal(t, synta.Filename, Filename{
		Segments: []Segment{
			{SegmentTypeIdentifier, &id_test, []Segment{
				{SegmentTypeOptional, nil, []Segment{
					{SegmentTypeIdentifier, &id_test, []Segment(nil)},
					{SegmentTypeOptional, nil, []Segment{
						{SegmentTypeIdentifier, &id_test, []Segment(nil)},
					}},
					{SegmentTypeIdentifier, &id_test, []Segment(nil)},
					{SegmentTypeOptional, nil, []Segment{
						{SegmentTypeIdentifier, &id_test, []Segment(nil)},
					}},
				}},
			}},
		},
		Extension: Identifier("test"),
	})
}

func TestParseSyntaWithNestedOptionalError(t *testing.T) {
	input := `; a test comment
; a second comment
test = a|b
> test(-test(-test)?)?)?.test`
	_, err := ParseSynta(input)
	assert.NotNil(t, err)
}
func TestParseSyntaWithNestedOptionalErrorBis(t *testing.T) {
	input := `; a test comment
; a second comment
test = a|b
> test(-test(-test)?.test`
	_, err := ParseSynta(input)
	assert.NotNil(t, err)
}
