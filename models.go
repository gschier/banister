package banister

import "strings"

type ModelSettings struct {
	Name        string
	VerboseName string
	DBTable     string
}

func (s *ModelSettings) Fix() {
	// Generate table name if it doesn't exist
	if s.DBTable == "" && strings.HasSuffix(s.Name, "s") {
		s.DBTable = DBName(s.Name)
	} else if s.DBTable == "" {
		s.DBTable = DBName(s.Name + "s")
	}
}

type Model interface {
	Settings() ModelSettings
	Fields() []Field
}
