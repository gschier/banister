package banister

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"strings"
)

type FilterGenerator struct {
	Model Model
	Field Field
	File  *jen.File
}

func NewFilterGenerator(file *jen.File, field Field, model Model) *FilterGenerator {
	return &FilterGenerator{Field: field, File: file, Model: model}
}

func (g *FilterGenerator) names() GeneratedFieldNames {
	return g.Field.Settings().Names(g.Model)
}

func (g *FilterGenerator) modelNames() GeneratedModelNames {
	return g.Model.Settings().Names()
}

func (g *FilterGenerator) goType() string {
	return fmt.Sprintf("%T", g.Field.EmptyDefault())
}

func (g *FilterGenerator) AddFilterMethod(name string, args *jen.Statement, filter *jen.Statement) {
	g.File.Func().Params(
		jen.Id("filter").Op("*").Id(g.names().FilterOptionStruct),
	).Id(name).Params(
		args,
	).Params(jen.Id(g.modelNames().QuerysetFilterArgStruct)).Block(
		jen.Return(
			jen.Id(g.modelNames().QuerysetFilterArgStruct).Values(jen.Dict{
				jen.Id("filter"): filter,
			}),
		),
	)
}

func (g *FilterGenerator) AddStructDef() {
	g.File.Type().Id(g.names().FilterOptionStruct).Struct(
		jen.Id("filters").Index().Qual("github.com/Masterminds/squirrel", "Sqlizer"),
	)
}

func (g *FilterGenerator) SqExpr(op Operation) *jen.Statement {
	switch op {
	case Exact:

	}

	return jen.Qual("github.com/Masterminds/squirrel", "Expr").Call()
}

func (g *FilterGenerator) TableDotColumn() string {
	return fmt.Sprintf(
		`"%s"."%s"`,
		strings.ReplaceAll(g.Model.Settings().DBTable, `"`, `\"`),
		strings.ReplaceAll(g.Field.Settings().DBColumn, `"`, `\"`),
	)
}

func (g *FilterGenerator) AddSimpleSquirrelFilter(name, sqName string) {
	g.AddFilterMethod(name,
		jen.Id("v").Id(g.goType()),
		jen.Op("&").Qual("github.com/Masterminds/squirrel", sqName).Values(jen.Dict{
			jen.Lit(g.TableDotColumn()): jen.Id("v"),
		}),
	)
}

func (g *FilterGenerator) Generate() {
	g.AddStructDef()
	for _, op := range g.Field.Operations() {
		switch op {
		case Exact:
			g.AddSimpleSquirrelFilter("Eq", "Eq")
		case Gt:
			g.AddSimpleSquirrelFilter("Gt", "Gt")
		case Gte:
			g.AddSimpleSquirrelFilter("Gte", "GtOrEq")
		case Lt:
			g.AddSimpleSquirrelFilter("Lt", "Lt")
		case Lte:
			g.AddSimpleSquirrelFilter("Lte", "LtOrEq")
		case Contains:
			g.AddSimpleSquirrelFilter("Contains", "Like")
		case IContains:
			g.AddSimpleSquirrelFilter("IContains", "ILike")
		}
	}

	if g.Field.Settings().Null {
		g.AddFilterMethod("Null",
			nil,
			jen.Op("&").Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Dict{
				jen.Lit(g.TableDotColumn()): jen.Nil(),
			}),
		)
	}
}
