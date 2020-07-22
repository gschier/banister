package banister

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"strings"
)

type FilterGenerator struct {
	Model Model
	File  *File
}

func NewFilterGenerator(file *File, model Model) *FilterGenerator {
	return &FilterGenerator{File: file, Model: model}
}

func (g *FilterGenerator) names(f Field) GeneratedFieldNames {
	return f.Settings().Names(g.Model)
}

func (g *FilterGenerator) modelNames() GeneratedModelNames {
	return g.Model.Settings().Names()
}

func (g *FilterGenerator) goType(f Field) string {
	return fmt.Sprintf("%T", f.EmptyDefault())
}

func (g *FilterGenerator) AddFilterMethod(f Field, name string, args *Statement, filter *Statement) {
	g.File.Func().Params(
		Id("filter").Op("*").Id(g.names(f).FilterOptionStruct),
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

func (g *FilterGenerator) AddStruct(f Field) {
	g.File.Type().Id(g.names(f).FilterOptionStruct).Struct(
		Id("filters").Index().Qual("github.com/Masterminds/squirrel", "Sqlizer"),
	)
}

func (g *FilterGenerator) SqExpr(op Operation) *Statement {
	switch op {
	case Exact:

	}

	return Qual("github.com/Masterminds/squirrel", "Expr").Call()
}

func (g *FilterGenerator) AddSimpleSquirrelFilter(f Field, name, sqName string) {
	// If type comes from package, we need to qualify it
	segments := strings.SplitN(g.goType(f), ".", 2)
	var defineGoType *Statement
	if len(segments) == 2 {
		defineGoType = Id("v").Qual(segments[0], segments[1])
	} else {
		defineGoType = Id("v").Id(segments[0])
	}

	g.AddFilterMethod(
		f,
		name,
		defineGoType,
		Op("&").Qual("github.com/Masterminds/squirrel", sqName).Values(Dict{
			Lit(g.names(f).QualifiedColumn): Id("v"),
		}),
	)
}

func (g *FilterGenerator) AddFilterOptionsStruct() {
	fields := make([]Code, 0)
	values := Dict{}
	for _, f := range g.Model.Fields() {
		fieldName := f.Settings().Name
		fieldType := f.Settings().Names(g.Model).FilterOptionStruct
		fields = append(fields, Id(fieldName).Id(fieldType))
		values[Id(fieldName)] = Id(fieldType).Values()
	}

	// Define options struct and instantiate an instance inline, so we
	// don't clutter the scope with unnecessary type.
	structName := g.Model.Settings().Names().QuerysetFilterOptionsStruct
	g.File.Type().Id(structName).Struct(fields...)
}

func (g *FilterGenerator) Generate() {
	for _, f := range g.Model.Fields() {
		g.AddStruct(f)
		for _, op := range f.Operations() {
			switch op {
			case Exact:
				g.AddSimpleSquirrelFilter(f, "Eq", "Eq")
			case Gt:
				g.AddSimpleSquirrelFilter(f, "Gt", "Gt")
			case Gte:
				g.AddSimpleSquirrelFilter(f, "Gte", "GtOrEq")
			case Lt:
				g.AddSimpleSquirrelFilter(f, "Lt", "Lt")
			case Lte:
				g.AddSimpleSquirrelFilter(f, "Lte", "LtOrEq")
			case Contains:
				g.AddSimpleSquirrelFilter(f, "Contains", "Like")
			case IContains:
				g.AddSimpleSquirrelFilter(f, "IContains", "ILike")
			}
		}

		if f.Settings().Null {
			g.AddFilterMethod(
				f,
				"Null",
				nil,
				Op("&").Qual("github.com/Masterminds/squirrel", "Eq").Values(Dict{
					Lit(g.names(f).QualifiedColumn): Nil(),
				}),
			)
		}
	}

	g.AddFilterOptionsStruct()
}
