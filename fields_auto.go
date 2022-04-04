package banister

func NewAutoField(name string) *AutoFieldBuilder {
	base := NewIntegerField(name).build()
	base.settings.PrimaryKey = true
	return &AutoFieldBuilder{
		field: &AutoField{base: &base},
	}
}

// AutoField is a database field used for integers
type AutoField struct {
	base *IntegerField
}

func (f AutoField) Type() FieldType {
	return Auto
}

func (f AutoField) RelType() FieldType {
	return f.base.Type()
}

func (f AutoField) EmptyDefault() interface{} {
	return f.base.EmptyDefault()
}

func (f AutoField) QueryOperators() map[QueryOperator]string {
	return f.base.QueryOperators()
}

func (f AutoField) Settings() FieldSettings {
	return f.base.Settings()
}

func (f AutoField) ProvideModels(_ Model, _ []Model) {
	// Nothing yet
}

type AutoFieldBuilder struct {
	field *AutoField
}

func (b *AutoFieldBuilder) Build() Field {
	return *b.field
}
