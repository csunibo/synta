package synta

// Clear returns a new Synta structure without any unused definitions
func Clear(synta Synta) (s Synta) {
	s.Filename = synta.Filename
	s.Definitions = map[Identifier]Definition{}
	s.Definitions[s.Filename.Extension] = synta.Definitions[s.Filename.Extension]
	clearSegments(synta, s, s.Filename.Segments)
	return
}

func clearSegments(synta Synta, s Synta, segments []Segment) {
	for _, segment := range segments {
		if segment.Kind == SegmentTypeIdentifier {
			s.Definitions[*segment.Value] = synta.Definitions[*segment.Value]
		} else if segment.Kind == SegmentTypeOptional {
			clearSegments(synta, s, segment.Subsegments)
		}
	}
}
