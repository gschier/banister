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

	ModelStruct      string
	ModelConstructor string
	ConfigStruct     string

	QuerysetFilterArgStruct  string
	QuerysetOrderByArgStruct string
	QuerysetSetterArgStruct  string

	QuerysetFilterOptionsStruct       string
	QuerysetOrderByOptionsStruct      string
	QuerysetOrderByOptionsConstructor string
	QuerysetSetterOptionsStruct       string
}

type GeneratedFieldNames struct {
	FilterOptionStruct             string
	OrderByOptionStruct            string
	SetterOptionStruct             string
	QualifiedColumn                string
	QuerysetOrderByDirectionStruct string
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
		ManagerStruct:       PrivateGoName(modelName + "Manager"),
		ManagerConstructor:  PrivateGoName("New" + modelName + "Manager"),
		QuerysetStruct:      PublicGoName(modelName + "Queryset"),
		QuerysetConstructor: PrivateGoName("New" + modelName + "Queryset"),

		ModelStruct:      PublicGoName(modelName),
		ModelConstructor: PublicGoName("New" + modelName),
		ConfigStruct:     PublicGoName(modelName + "Config"),

		QuerysetFilterOptionsStruct:       PrivateGoName(modelName + "Filters"),
		QuerysetOrderByOptionsStruct:      PrivateGoName(modelName + "Orders"),
		QuerysetOrderByOptionsConstructor: PrivateGoName("new" + modelName + "Orders"),
		QuerysetSetterOptionsStruct:       PrivateGoName(modelName + "Setters"),

		QuerysetFilterArgStruct:  PrivateGoName(modelName + "FilterArg"),
		QuerysetOrderByArgStruct: PrivateGoName(modelName + "OrderByArg"),
		QuerysetSetterArgStruct:  PrivateGoName(modelName + "SetterArg"),
	}
}

func NamesForField(modelSettings ModelSettings, fieldSettings FieldSettings) GeneratedFieldNames {
	modelName := modelSettings.Name
	fieldName := fieldSettings.Name

	return GeneratedFieldNames{
		FilterOptionStruct:             PrivateGoName(modelName + fieldName + "Filter"),
		OrderByOptionStruct:            PrivateGoName(modelName + fieldName + "OrderBy"),
		SetterOptionStruct:             PrivateGoName(modelName + fieldName + "Setter"),
		QuerysetOrderByDirectionStruct: PrivateGoName(modelName + fieldName + "OrderByArg"),

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
