package banister

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"strings"
)

type FilterGenerator struct {
	Model Model
	Field Field
	File  *File
}

func NewFilterGenerator(file *File, field Field, model Model) *FilterGenerator {
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

func (g *FilterGenerator) AddFilterMethod(name string, args *Statement, filter *Statement) {
	g.File.Func().Params(
		Id("filter").Op("*").Id(g.names().FilterOptionStruct),
	).Id(name).Params(
		args,
	).Params(Id(g.modelNames().QuerysetFilterArgStruct)).Block(
		Return(
			Id(g.modelNames().QuerysetFilterArgStruct).Values(Dict{
				Id("filter"): filter,
			}),
		),
	)
}

func (g *FilterGenerator) AddStruct() {
	g.File.Type().Id(g.names().FilterOptionStruct).Struct(
		Id("filters").Index().Qual("github.com/Masterminds/squirrel", "Sqlizer"),
	)
}

func (g *FilterGenerator) SqExpr(op Operation) *Statement {
	switch op {
	case Exact:

	}

	return Qual("github.com/Masterminds/squirrel", "Expr").Call()
}

func (g *FilterGenerator) TableDotColumn() string {
	return fmt.Sprintf(
		`"%s"."%s"`,
		strings.ReplaceAll(g.Model.Settings().DBTable, `"`, `\"`),
		strings.ReplaceAll(g.Field.Settings().DBColumn, `"`, `\"`),
	)
}

func (g *FilterGenerator) AddSimpleSquirrelFilter(name, sqName string) {
	// If type comes from package, we need to qualify it
	segments := strings.SplitN(g.goType(), ".", 2)
	var defineGoType *Statement
	if len(segments) == 2 {
		defineGoType = Id("v").Qual(segments[0], segments[1])
	} else {
		defineGoType = Id("v").Id(segments[0])
	}

	g.AddFilterMethod(name,
		defineGoType,
		Op("&").Qual("github.com/Masterminds/squirrel", sqName).Values(Dict{
			Lit(g.TableDotColumn()): Id("v"),
		}),
	)
}

func (g *FilterGenerator) Generate() {
	g.AddStruct()
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
			Op("&").Qual("github.com/Masterminds/squirrel", "Eq").Values(Dict{
				Lit(g.TableDotColumn()): Nil(),
			}),
		)
	}
}
