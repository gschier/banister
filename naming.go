package banister

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/takuoki/gocase"
	"strings"
)

type GeneratedModelNames struct {
	ManagerAccessor     string
	ManagerStruct       string
	ManagerConstructor  string
	QuerysetStruct      string
	QuerysetConstructor string

	ModelStruct  string
	ConfigStruct string

	QuerysetFilterArgStruct  string
	QuerysetOrderByArgStruct string
	QuerysetSetterArgStruct  string

	QuerysetFilterOptionsStruct  string
	QuerysetOrderByOptionsStruct string
	QuerysetSetterOptionsStruct  string

	FilterOptionsVar  string
	OrderByOptionsVar string
	SetterOptionsVar  string
}

type GeneratedFieldNames struct {
	FilterOptionStruct  string
	OrderByOptionStruct string
	SetterOptionStruct  string
	QualifiedColumn     string
}

var globalNames = struct {
	StoreStruct       string
	StoreConfigStruct string
	StoreConstructor  string

	HookPostInsert string
	HookPreInsert  string
	HookPostUpdate string
	HookPreUpdate  string
	HookPostDelete string
	HookPreDelete  string
}{
	StoreStruct:       PublicGoName("Store"),
	StoreConfigStruct: PublicGoName("StoreConfig"),
	StoreConstructor:  PublicGoName("NewStore"),

	HookPostInsert: "HookPostInsert",
	HookPreInsert:  "HookPreInsert",
	HookPostUpdate: "HookPostUpdate",
	HookPreUpdate:  "HookPreUpdate",
	HookPostDelete: "HookPostDelete",
	HookPreDelete:  "HookPreDelete",
}

func NamesForModel(modelName string) GeneratedModelNames {
	return GeneratedModelNames{
		ManagerAccessor:     PublicGoName(modelName + "s"),
		ManagerStruct:       PublicGoName(modelName + "Manager"),
		ManagerConstructor:  PublicGoName("New" + modelName + "Manager"),
		QuerysetStruct:      PublicGoName(modelName + "Queryset"),
		QuerysetConstructor: PublicGoName("New" + modelName + "Queryset"),

		ModelStruct:  PublicGoName(modelName),
		ConfigStruct: PublicGoName(modelName + "Config"),

		QuerysetFilterOptionsStruct:  PublicGoName(modelName + "Filters"),
		QuerysetOrderByOptionsStruct: PublicGoName(modelName + "OrderBys"),
		QuerysetSetterOptionsStruct:  PublicGoName(modelName + "Setters"),

		QuerysetFilterArgStruct:  PrivateGoName(modelName + "FilterArg"),
		QuerysetOrderByArgStruct: PrivateGoName(modelName + "OrderByArg"),
		QuerysetSetterArgStruct:  PrivateGoName(modelName + "SetterArg"),

		FilterOptionsVar:  PublicGoName("Where" + modelName),
		OrderByOptionsVar: PublicGoName("OrderBy" + modelName),
		SetterOptionsVar:  PublicGoName("Set" + modelName),
	}
}

func NamesForField(modelSettings ModelSettings, fieldSettings FieldSettings) GeneratedFieldNames {
	modelName := modelSettings.Name
	fieldName := fieldSettings.Name

	return GeneratedFieldNames{
		FilterOptionStruct:  PublicGoName(modelName + fieldName + "Filter"),
		OrderByOptionStruct: PublicGoName(modelName + fieldName + "OrderBy"),
		SetterOptionStruct:  PublicGoName(modelName + fieldName + "Setter"),

		QualifiedColumn: fmt.Sprintf(
			`"%s"."%s"`,
			strings.ReplaceAll(modelSettings.DBTable, `"`, `\"`),
			strings.ReplaceAll(fieldSettings.DBColumn, `"`, `\"`),
		),
	}
}

func PublicGoName(rawName string) string {
	return gocase.To(strcase.ToCamel(rawName))
}

func PrivateGoName(rawName string) string {
	return gocase.To(strcase.ToLowerCamel(rawName))
}

func DBName(rawName string) string {
	return strcase.ToSnake(rawName)
}

func JSONName(rawName string) string {
	// NOTE: Converting to snake first removes consecutive UPPER
	return strcase.ToLowerCamel(strcase.ToSnake(rawName))
}
