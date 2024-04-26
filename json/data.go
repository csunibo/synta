package json

// A Definition is a named regexp along with comments
// to clarify the regexp's purpose
// It corresponds to the <commdef> BNF definition
type Definition struct {
	Comments []string `json:"comments"`
	Regexp   string   `json:"regexp"`
}

// A Segment is a section of the main filename
// It corresponds to the <segment> BNF definition
type Segment struct {
	Kind        uint      `json:"kind"`
	Value       string    `json:"value"`
	Subsegments []Segment `json:"subsegments"`
}

// Filename represents the flename defintion, made up
// of a series of segments and a file extension
type Filename struct {
	Segments  []Segment `json:"segments"`
	Extension string    `json:"extension"`
}

// Synta represents the contents of a Synta file
// It corresponds to the <language> BNF definition
// The last segment of the Filename is the extension
type Synta struct {
	Definitions map[string]Definition `json:"definitions"`
	Filename    Filename              `json:"filename"`
}
