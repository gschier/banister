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

func (s ModelSettings) PluralName() string {
	if strings.HasSuffix(s.Name, "s") {
		return s.Name
	}

	return s.Name + "s"
}

func (s *ModelSettings) FillDefaults() {
	// Generate table name if it doesn't exist
	if s.DBTable == "" {
		s.DBTable = DBName(s.PluralName())
	}

	if s.VerboseName == "" {
		s.VerboseName = s.Name
	}
}

type Model interface {
	Settings() ModelSettings
	Fields() []Field
	ProvideModels(models ...Model)
}

type model struct {
	settings      *ModelSettings
	fieldBuilders []FieldBuilder
	fields        []Field
}

func NewModel(name string, field ...FieldBuilder) Model {
	fields := make([]Field, len(field))
	for i, b := range field {
		fields[i] = b.Build()
	}

	settings := &ModelSettings{Name: name}
	settings.FillDefaults()

	return model{
		settings:      settings,
		fieldBuilders: field,
		fields:        fields,
	}
}

func (m model) Settings() ModelSettings {
	return *m.settings
}

// ProvideModels gives the model an opportunity to setup things that may
// require knowledge of the rest of the models
func (m model) ProvideModels(models ...Model) {
	for _, f := range m.fields {
		f.ProvideModels(m, models)
	}
}

func (m model) Fields() []Field {
	if m.fields == nil {
		panic(m.settings.Name + " model has not been initialized")
	}
	return m.fields
}

func PrimaryKeyField(m Model) Field {
	for _, f := range m.Fields() {
		if f.Settings().PrimaryKey {
			return f
		}
	}

	panic("no primary key field for model: " + m.Settings().Name)
}
