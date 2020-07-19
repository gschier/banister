package banister

import (
	"github.com/iancoleman/strcase"
	"github.com/takuoki/gocase"
)

type GeneratedModelNames struct {
	StoreStruct               string
	ManagerStruct             string
	QuerysetStruct            string
	QuerysetStructConstructor string

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
}

func NamesForModel(modelName string) GeneratedModelNames {
	return GeneratedModelNames{
		StoreStruct:               PublicGoName("Store"),
		ManagerStruct:             PublicGoName(modelName + "Manager"),
		QuerysetStruct:            PublicGoName(modelName + "Queryset"),
		QuerysetStructConstructor: PublicGoName("New" + modelName + "Queryset"),

		ModelStruct:  PublicGoName(modelName),
		ConfigStruct: PublicGoName(modelName + "Config"),

		QuerysetFilterOptionsStruct:  PublicGoName(modelName + "QuerysetFilterOptions"),
		QuerysetOrderByOptionsStruct: PublicGoName(modelName + "QuerysetOrderByOptions"),
		QuerysetSetterOptionsStruct:  PublicGoName(modelName + "QuerysetSetterOptions"),

		QuerysetFilterArgStruct:  PrivateGoName(modelName + "QuerysetFilterArg"),
		QuerysetOrderByArgStruct: PrivateGoName(modelName + "QuerysetOrderByArg"),
		QuerysetSetterArgStruct:  PrivateGoName(modelName + "QuerysetSetterArg"),

		FilterOptionsVar:  PublicGoName("Where" + modelName),
		OrderByOptionsVar: PublicGoName("OrderBy" + modelName),
		SetterOptionsVar:  PublicGoName("Set" + modelName),
	}
}

func NamesForField(modelName, fieldName string) GeneratedFieldNames {
	return GeneratedFieldNames{
		FilterOptionStruct:  PublicGoName(modelName + fieldName + "Filter"),
		OrderByOptionStruct: PublicGoName(modelName + fieldName + "OrderBy"),
		SetterOptionStruct:  PublicGoName(modelName + fieldName + "Setter"),
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
