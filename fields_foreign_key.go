package banister

import "strings"

func NewForeignKeyField(to string) *ForeignKeyFieldBuilder {
	return &ForeignKeyFieldBuilder{
		field: &ForeignKeyField{
			to: to,
			settings: &FieldSettings{
				Name:      PublicGoName(to + "ID"),
				ManyToOne: true,
				Rel:       &Rel{},
			},
		},
	}
}

// ForeignKeyField
type ForeignKeyField struct {
	settings *FieldSettings
	to       string
}

func (f ForeignKeyField) Settings() FieldSettings {
	return *f.settings
}

func (f ForeignKeyField) Type() FieldType {
	return f.settings.Rel.ToField.RelType()
}

func (f ForeignKeyField) RelType() FieldType {
	return f.Type()
}

func (f ForeignKeyField) EmptyDefault() interface{} {
	// If we are an integer, set the default to -1, because ID 0 is likely
	// to exist
	// TODO: Force FKs to be set explicitly by the user (maybe make it a pointer?)
	if f.settings.Rel.ToField.RelType() == Integer {
		return int64(-1)
	}

	return f.settings.Rel.ToField.EmptyDefault()
}

func (f ForeignKeyField) QueryOperators() map[QueryOperator]string {
	return f.settings.Rel.ToField.QueryOperators()
}

func (f ForeignKeyField) ProvideModels(parent Model, models []Model) {
	var toModel Model
	for _, m := range models {
		if strings.ToLower(m.Settings().Name) == strings.ToLower(f.to) {
			toModel = m
			break
		}
	}

	if toModel == nil {
		panic("Failed to find related model " + f.to)
	}

	toField := PrimaryKeyField(toModel)

	if toField.Type() == Char {
		f.settings.MaxLength = toField.Settings().MaxLength
	}

	f.settings.Rel.To = toModel
	f.settings.Rel.ToField = toField
	f.settings.Rel.RelatedName = parent.Settings().PluralName()

	// TODO: Handle these
	//	f.settings.Rel.RelatedQueryName= ""
}

// ForeignKeyFieldBuilder
type ForeignKeyFieldBuilder struct {
	field *ForeignKeyField
}

func (f *ForeignKeyFieldBuilder) Build() Field {
	f.field.settings.Fix()
	return *f.field
}

func (f *ForeignKeyFieldBuilder) Null() *ForeignKeyFieldBuilder {
	f.field.settings.Null = true
	return f
}

func (f *ForeignKeyFieldBuilder) OnDelete(option OnDeleteOption) *ForeignKeyFieldBuilder {
	f.field.settings.Rel.OnDelete = option.OnDeleteValue()
	return f
}

type OnDeleteOption interface {
	OnDeleteValue() string
}

type onDeleteValue struct {
	value string
}

func (d onDeleteValue) OnDeleteValue() string {
	return d.value
}

var (
	OnDeleteSetNull    = onDeleteValue{value: "SET NULL"}
	OnDeleteSetDefault = onDeleteValue{value: "SET DEFAULT"}
	OnDeleteRestrict   = onDeleteValue{value: "RESTRICT"}
	OnDeleteNoAction   = onDeleteValue{value: "NO ACTION"}
	OnDeleteCascade    = onDeleteValue{value: "CASCADE"}
)
