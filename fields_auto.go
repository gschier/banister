package banister

func NewAutoField(name string) *AutoFieldBuilder {
	settings := NewIntegerField(name).Build().settings
	settings.PrimaryKey = true
	return &AutoFieldBuilder{
		field: &AutoField{settings: settings},
	}
}

// IntegerField is a database field used for integers
type AutoField struct {
	settings *FieldSettings
}

func (f AutoField) Type() FieldType {
	return Auto
}

func (f AutoField) RelType() FieldType {
	return Integer
}

func (f AutoField) EmptyDefault() interface{} {
	return int64(0)
}

func (f AutoField) Operations() map[Operation]string {
	return map[Operation]string{
		Exact: "Eq",
		Lt:    "Lt",
		Lte:   "Lte",
		Gt:    "Gt",
		Gte:   "Gte",
	}
}

func (f AutoField) Settings() FieldSettings {
	f.settings.Fix()
	return *f.settings
}

// AutoFieldBuilder
type AutoFieldBuilder struct {
	field *AutoField
}

func (b *AutoFieldBuilder) Build() AutoField {
	return *b.field
}
