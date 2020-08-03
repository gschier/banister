package banister

func NewTextArrayField(name string) *TextArrayFieldBuilder {
	return &TextArrayFieldBuilder{
		field: &TextArrayField{
			settings: &FieldSettings{Name: name},
		},
	}
}

// TextArrayField
type TextArrayField struct {
	settings *FieldSettings
}

func (f TextArrayField) Settings() FieldSettings {
	return *f.settings
}

func (f TextArrayField) Type() FieldType {
	return TextArray
}

func (f TextArrayField) RelType() FieldType {
	return TextArray
}

func (f TextArrayField) EmptyDefault() interface{} {
	return []string{}
}

func (f TextArrayField) QueryOperators() map[QueryOperator]string {
	return map[QueryOperator]string{
		ArrayContains: "Contains",
	}
}

func (f TextArrayField) ProvideModels(_ Model, _ []Model) {
	// Nothing yet
}

// TextArrayFieldBuilder
type TextArrayFieldBuilder struct {
	field *TextArrayField
}

func (f *TextArrayFieldBuilder) Build() Field {
	return f.build()
}

func (f *TextArrayFieldBuilder) build() TextArrayField {
	f.field.settings.Fix()
	return *f.field
}

func (f *TextArrayFieldBuilder) Default(s string) *TextArrayFieldBuilder {
	f.field.settings.Default = NewDefaultValue(s)
	return f
}

func (f *TextArrayFieldBuilder) Null() *TextArrayFieldBuilder {
	f.field.settings.Null = true
	return f
}

func (f *TextArrayFieldBuilder) Unique() *TextArrayFieldBuilder {
	f.field.settings.Unique = true
	return f
}

func (f *TextArrayFieldBuilder) Size(size uint) *TextArrayFieldBuilder {
	f.field.settings.Size = &size
	return f
}
