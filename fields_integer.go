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

func (f IntegerField) Operations() map[Operation]string {
	return map[Operation]string{
		Exact: "Eq",
		Lt:    "Lt",
		Lte:   "Lte",
		Gt:    "Gt",
		Gte:   "Gte",
	}
}

func (f IntegerField) ProvideModels(_ Model, _ []Model) {
	// Nothing yet
}

// IntegerFieldBuilder
type IntegerFieldBuilder struct {
	field *IntegerField
}

func (f *IntegerFieldBuilder) build() IntegerField {
	f.field.settings.Fix()
	return *f.field
}

func (f *IntegerFieldBuilder) Build() Field {
	return f.build()
}

func (f *IntegerFieldBuilder) Default(i int64) *IntegerFieldBuilder {
	f.field.settings.Default = NewDefaultValue(i)
	return f
}

func (f *IntegerFieldBuilder) Null() *IntegerFieldBuilder {
	f.field.settings.Null = true
	return f
}
