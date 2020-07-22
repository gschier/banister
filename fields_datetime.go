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

func (f DateTimeField) Operations() map[Operation]string {
	return map[Operation]string{
		Exact: "Eq",
		Lt:    "Before",
		Gt:    "After",
	}
}

// DateTimeFieldBuilder
type DateTimeFieldBuilder struct {
	field *DateTimeField
}

func (f *DateTimeFieldBuilder) Build() DateTimeField {
	f.field.settings.Fix()
	return *f.field
}

func (f *DateTimeFieldBuilder) Hidden() *DateTimeFieldBuilder {
	f.field.settings.Hidden = true
	return f
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
