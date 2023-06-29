package synta

// Clear returns a new Synta structure without any unused definitions
func Clear(synta Synta) (s Synta) {
	s.Filename = synta.Filename
	s.Definitions = map[Identifier]Definition{}
	s.Definitions[s.Filename.Extension] = synta.Definitions[s.Filename.Extension]
	for _, segment := range s.Filename.Segments {
		s.Definitions[segment.Identifier] = synta.Definitions[segment.Identifier]
	}
	return
}
