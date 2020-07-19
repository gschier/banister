package banister

import (
	"fmt"
	"github.com/dave/jennifer/jen"
)

type FilterGenerator struct {
	Field Field
	File  *jen.File
}

func NewFilterGenerator(file *jen.File, field Field) *FilterGenerator {
	return &FilterGenerator{Field: field, File: file}
}

func (g *FilterGenerator) StructName() string {
	return g.Field.Settings().Name + "Filter"
}

func (g *FilterGenerator) GoType() string {
	return fmt.Sprintf("%T", g.Field.EmptyDefault())
}

func (g *FilterGenerator) AddFilterMethod(name, sqlizerName string, args *jen.Statement, lit *jen.Statement) {
	g.File.Comment(name + " does a thing")
	g.File.Func().Params(
		jen.Id("filter").Op("*").Id(g.StructName()),
	).Id(name).Params(
		args,
	).Params(jen.Op("*").Id(g.StructName())).Block(
		jen.Id("filter").Dot("filters").Op("=").Id("append").Params(
			jen.Id("filter"),
			jen.Op("&").Qual("github.com/Masterminds/squirrel", sqlizerName).Values(
				jen.Dict{jen.Lit(g.Field.Settings().DBColumn): lit},
			),
		),
		jen.Return(jen.Id("filter")),
	)
}

func (g *FilterGenerator) AddStructDef() {
	g.File.Type().Id(g.StructName()).Struct(
		jen.Id("filters").Index().Qual("github.com/Masterminds/squirrel", "Sqlizer"),
	)
}

func (g *FilterGenerator) Generate() {
	g.AddStructDef()
	g.AddFilterMethod("Eq", "Eq", jen.Id("v").String(), jen.Lit("v"))
	g.AddFilterMethod("Gt", "Gt", jen.Id("v").String(), jen.Lit("v"))
	g.AddFilterMethod("Gte", "GtOrEq", jen.Id("v").String(), jen.Lit("v"))
	g.AddFilterMethod("Lt", "Lt", jen.Id("v").String(), jen.Lit("v"))
	g.AddFilterMethod("Lte", "LtOrEq", jen.Id("v").String(), jen.Lit("v"))

	if g.Field.Settings().Null {
		g.AddFilterMethod("Null", "Eq", nil, jen.Nil())
	}
}
