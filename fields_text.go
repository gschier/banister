package banister

func NewTextField(name string) *TextFieldBuilder {
	return &TextFieldBuilder{
		field: &TextField{
			settings: &FieldSettings{Name: name},
		},
	}
}

// TextField
type TextField struct {
	settings *FieldSettings
}

func (f TextField) Settings() FieldSettings {
	return *f.settings
}

func (f TextField) Type() FieldType {
	return Text
}

func (f TextField) RelType() FieldType {
	return Text
}

func (f TextField) EmptyDefault() interface{} {
	return ""
}

// TextFieldBuilder
type TextFieldBuilder struct {
	field *TextField
}

func (f *TextFieldBuilder) Build() TextField {
	f.field.settings.Fix()
	return *f.field
}

func (f *TextFieldBuilder) Hidden() *TextFieldBuilder {
	f.field.settings.Hidden = true
	return f
}

func (f *TextFieldBuilder) Default(s string) *TextFieldBuilder {
	f.field.settings.Default = NewDefaultValue(s)
	return f
}

func (f *TextFieldBuilder) Null() *TextFieldBuilder {
	f.field.settings.Null = true
	return f
}

func (f *TextFieldBuilder) Unique() *TextFieldBuilder {
	f.field.settings.Unique = true
	return f
}
