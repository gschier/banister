package banister

func NewBooleanField(name string) *BooleanFieldBuilder {
	return &BooleanFieldBuilder{
		field: &BooleanField{
			settings: &FieldSettings{Name: name},
		},
	}
}

// BooleanField
type BooleanField struct {
	settings *FieldSettings
}

func (f BooleanField) Settings() FieldSettings {
	return *f.settings
}

func (f BooleanField) Type() FieldType {
	return Boolean
}

func (f BooleanField) RelType() FieldType {
	return Boolean
}

func (f BooleanField) EmptyDefault() interface{} {
	return false
}

func (f BooleanField) Operations() []Operation {
	return []Operation{Exact}
}

// BooleanFieldBuilder
type BooleanFieldBuilder struct {
	field *BooleanField
}

func (f *BooleanFieldBuilder) Build() BooleanField {
	f.field.settings.Fix()
	return *f.field
}

func (f *BooleanFieldBuilder) Hidden() *BooleanFieldBuilder {
	f.field.settings.Hidden = true
	return f
}

func (f *BooleanFieldBuilder) Default(s bool) *BooleanFieldBuilder {
	f.field.settings.Default = NewDefaultValue(s)
	return f
}

func (f *BooleanFieldBuilder) Null() *BooleanFieldBuilder {
	f.field.settings.Null = true
	return f
}
