package banister

import (
	"github.com/iancoleman/strcase"
	"github.com/takuoki/gocase"
)

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
