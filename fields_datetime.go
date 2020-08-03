package banister

import "time"

func NewDateTimeField(name string) *DateTimeFieldBuilder {
	return &DateTimeFieldBuilder{
		field: &DateTimeField{
			settings: &FieldSettings{Name: name},
		},
	}
}

// DateTimeField
type DateTimeField struct {
	settings *FieldSettings
}

func (f DateTimeField) Settings() FieldSettings {
	return *f.settings
}

func (f DateTimeField) Type() FieldType {
	return DateTime
}

func (f DateTimeField) RelType() FieldType {
	return DateTime
}

func (f DateTimeField) EmptyDefault() interface{} {
	return time.Time{}
}

func (f DateTimeField) QueryOperators() map[QueryOperator]string {
	return map[QueryOperator]string{
		Exact:    "Eq",
		NotExact: "NotEq",
		Lt:       "Before",
		Gt:       "After",
	}
}

func (f DateTimeField) ProvideModels(_ Model, _ []Model) {
	// Nothing yet
}

// DateTimeFieldBuilder
type DateTimeFieldBuilder struct {
	field *DateTimeField
}

func (f *DateTimeFieldBuilder) Build() Field {
	f.field.settings.Fix()
	return *f.field
}

func (f *DateTimeFieldBuilder) Default(s time.Time) *DateTimeFieldBuilder {
	f.field.settings.Default = NewDefaultValue(s)
	return f
}

func (f *DateTimeFieldBuilder) Null() *DateTimeFieldBuilder {
	f.field.settings.Null = true
	return f
}

func (f *DateTimeFieldBuilder) Unique() *DateTimeFieldBuilder {
	f.field.settings.Unique = true
	return f
}
