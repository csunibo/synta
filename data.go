package synta

import (
	"regexp"
)

// A regexp which describes an identifier
var IdentifierRegexp = regexp.MustCompile("[a-z]+")

// An Identifier is a lowercase alphabetical string.
// It corresponds to the <id> BNF definition
type Identifier string

// A Definition is a named regexp along with comments
// to clarify the regexp's purpose
// It corresponds to the <commdef> BNF definition
type Definition struct {
	Comments []string
	Regexp   *regexp.Regexp
}

type SegmentType uint

const (
	SegmentTypeIdentifier = iota
	SegmentTypeOptional
)

// A Segment is a section of the main filename
// It corresponds to the <segment> BNF definition
type Segment struct {
	kind        SegmentType
	value       *Identifier
	subsegments []Segment
}

// Filename represents the flename defintion, made up
// of a series of segments and a file extension
type Filename struct {
	Segments  []Segment
	Extension Identifier
}

// Synta represents the contents of a Synta file
// It corresponds to the <language> BNF definition
// The last segment of the Filename is the extension
type Synta struct {
	Definitions map[Identifier]Definition
	Filename    Filename
}
