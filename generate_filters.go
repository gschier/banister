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
	createFilter := Return(
		Id(g.modelNames().QuerysetFilterArgStruct).Values(Dict{
			Id("filter"): filter,
		}),
	)

	g.AddMethodField(f,
		name,
		[]Code{args},
		[]Code{createFilter},
		[]Code{Id("filter").Id(g.modelNames().QuerysetFilterArgStruct)},
	)
}

func (g *FilterGenerator) AddStruct(f Field) {
	g.File.Type().Id(g.names(f).FilterOptionStruct).Struct()
}

func (g *FilterGenerator) AddMethodField(f Field, name string, args, block, returns []Code) {
	receiver := Id(g.names(f).FilterOptionStruct)
	g.File.Func().Params(receiver).Id(name).
		Params(args...).
		Params(returns...).
		Block(block...)
}

func (g *FilterGenerator) AddMethod(name string, args, block, returns []Code) {
	receiver := Id(g.modelNames().QuerysetFilterOptionsStruct)
	g.File.Func().Params(receiver).Id(name).
		Params(args...).
		Params(returns...).
		Block(block...)
}

func (g *FilterGenerator) AddFilterExprMethod(f Field, name string, op Operation) {
	exprStr, ok := __backend.FilterOperations()[op]
	if !ok {
		panic("Unsupported filter operation " + name)
	}

	// If type comes from package, we need to qualify it
	segments := strings.SplitN(g.goType(f), ".", 2)
	var defineGoType *Statement
	if len(segments) == 2 {
		defineGoType = Id("v").Qual(segments[0], segments[1])
	} else {
		defineGoType = Id("v").Id(segments[0])
	}

	Expr := Qual("github.com/Masterminds/squirrel", "Expr")
	filterDef := Expr.Call(Lit(g.names(f).QualifiedColumn+" "+exprStr), Id("v"))
	g.AddFilterMethod(f, name, defineGoType, filterDef)
}

func (g *FilterGenerator) AddAndOrMethods() {
	// TODO: Implement And() and Or() as generic functions when/if generics
	//   are added to Go.
	for fnName, sqType := range map[string]string{"And": "And", "Or": "Or"} {
		sqType := Qual("github.com/Masterminds/squirrel", sqType)
		sqDef := Id("q").Op(":=").Add(sqType).Values()
		joinDef := Id("j").Op(":=").Make(Index().String(), Lit(0))

		mapFilters := For(
			Op("_").Op(",").Id("f").Op(":=").Range().Id("filter"),
		).Block(
			Op("q").Op("=").Append(Id("q"), Id("f").Dot("filter")),
			Op("j").Op("=").Append(Id("j"), Id("f").Dot("joins").Op("...")),
		)

		filter := Id(g.modelNames().QuerysetFilterArgStruct).Values(Dict{
			Id("filter"): Id("q"),
			Id("joins"):  Id("j"),
		})

		g.File.Comment(fnName + " combines multiple filters into one")
		g.AddMethod(fnName,
			[]Code{Id("filter").Op("...").Id(g.modelNames().QuerysetFilterArgStruct)},
			[]Code{
				sqDef,
				joinDef.Line(),
				mapFilters.Line(),
				Return(filter),
			},
			[]Code{Id(g.modelNames().QuerysetFilterArgStruct)},
		)
	}
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

func (g *FilterGenerator) AddOperationFilters(f Field) {
	for op, name := range f.Operations() {
		g.AddFilterExprMethod(f, name, op)
	}
}

func (g *FilterGenerator) AddNullFilterMaybe(f Field) {
	if !f.Settings().Null {
		return
	}

	eq := Op("&").Qual("github.com/Masterminds/squirrel", "Eq")
	filter := eq.Values(Dict{Lit(g.names(f).QualifiedColumn): Nil()})
	g.AddFilterMethod(f, "IsNull", nil, filter)
}

func (g *FilterGenerator) Generate() {
	for _, f := range g.Model.Fields() {
		g.AddStruct(f)
		g.AddOperationFilters(f)
		g.AddNullFilterMaybe(f)
	}

	g.AddFilterOptionsStruct()
	g.AddAndOrMethods()
}
