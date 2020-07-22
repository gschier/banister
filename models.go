package banister

import "strings"

type ModelSettings struct {
	Name        string
	VerboseName string
	DBTable     string
}

func (s ModelSettings) Names() GeneratedModelNames {
	return NamesForModel(s.Name)
}

func (s *ModelSettings) FillDefaults() {
	// Generate table name if it doesn't exist
	if s.DBTable == "" && strings.HasSuffix(s.Name, "s") {
		s.DBTable = DBName(s.Name)
	} else if s.DBTable == "" {
		s.DBTable = DBName(s.Name + "s")
	}

	if s.VerboseName == "" {
		s.VerboseName = s.Name
	}
}

type Model interface {
	Settings() ModelSettings
	Fields() []Field
}

func PrimaryKeyField(m Model) Field {
	for _, f := range m.Fields() {
		if f.Settings().PrimaryKey {
			return f
		}
	}

	panic("no primary key field for model: " + m.Settings().Name)
}
