package banister

func NewCharField(name string, maxLength uint) *CharFieldBuilder {
	base := NewTextField(name).build()
	base.settings.MaxLength = &maxLength
	return &CharFieldBuilder{
		field: &CharField{
			operations: base.QueryOperators(),
			settings:   base.settings,
		},
	}
}

// CharField
type CharField struct {
	operations map[QueryOperator]string
	settings   *FieldSettings
}

func (f CharField) Settings() FieldSettings {
	return *f.settings
}

func (f CharField) Type() FieldType {
	return Char
}

func (f CharField) RelType() FieldType {
	return Char
}

func (f CharField) EmptyDefault() interface{} {
	return ""
}

func (f CharField) QueryOperators() map[QueryOperator]string {
	return f.operations
}

func (f CharField) ProvideModels(_ Model, _ []Model) {
	// Nothing yet
}

// CharFieldBuilder
type CharFieldBuilder struct {
	field *CharField
}

func (f *CharFieldBuilder) Build() Field {
	f.field.settings.Fix()
	return *f.field
}

func (f *CharFieldBuilder) Default(s string) *CharFieldBuilder {
	f.field.settings.Default = NewDefaultValue(s)
	return f
}

func (f *CharFieldBuilder) Null() *CharFieldBuilder {
	f.field.settings.Null = true
	return f
}

func (f *CharFieldBuilder) Unique() *CharFieldBuilder {
	f.field.settings.Unique = true
	return f
}

func (f *CharFieldBuilder) PrimaryKey() *CharFieldBuilder {
	f.field.settings.PrimaryKey = true
	return f
}
