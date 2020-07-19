package banister

func NewIntegerField(name string) *IntegerFieldBuilder {
	return &IntegerFieldBuilder{
		field: &IntegerField{
			settings: &FieldSettings{Name: name},
		},
	}
}

// IntegerField
type IntegerField struct {
	settings *FieldSettings
}

func (f IntegerField) Settings() FieldSettings {
	return *f.settings
}

func (f IntegerField) Type() FieldType {
	return Integer
}

func (f IntegerField) RelType() FieldType {
	return Integer
}

func (f IntegerField) EmptyDefault() interface{} {
	return int64(0)
}

func (f IntegerField) Operations() []Operation {
	return []Operation{Exact, Lt, Lte, Gt, Gte}
}

// IntegerFieldBuilder
type IntegerFieldBuilder struct {
	field *IntegerField
}

func (f *IntegerFieldBuilder) Build() IntegerField {
	f.field.settings.Fix()
	return *f.field
}

func (f *IntegerFieldBuilder) Hidden(b bool) *IntegerFieldBuilder {
	f.field.settings.Hidden = b
	return f
}

func (f *IntegerFieldBuilder) Default(i int64) *IntegerFieldBuilder {
	f.field.settings.Default = NewDefaultValue(i)
	return f
}

func (f *IntegerFieldBuilder) Null() *IntegerFieldBuilder {
	f.field.settings.Null = true
	return f
}
