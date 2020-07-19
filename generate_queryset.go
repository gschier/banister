package banister

import "github.com/dave/jennifer/jen"

type QuerysetGenerator struct {
	File  *jen.File
	Model Model
}

func NewQuerysetGenerator(file *jen.File, model Model) *QuerysetGenerator {
	return &QuerysetGenerator{File: file, Model: model}
}

func (g *QuerysetGenerator) names() GeneratedModelNames {
	return g.Model.Settings().Names()
}

func (g *QuerysetGenerator) AddFilterArgsStruct() {
	g.File.Type().Id(g.names().QuerysetFilterArgStruct).Struct(
		jen.Id("filter").Qual("github.com/Masterminds/squirrel", "Sqlizer"),
		jen.Id("joins").Index().String(),
	)
}

func (g *QuerysetGenerator) AddOrderByArgsStruct() {
	g.File.Type().Id(g.names().QuerysetOrderByArgStruct).Struct(
		jen.Id("field").String(),
		jen.Id("order").String(),
		jen.Id("join").String(),
	)
}

func (g *QuerysetGenerator) AddSetterArgsStruct() {
	g.File.Type().Id(g.names().QuerysetSetterArgStruct).Struct(
		jen.Id("field").String(),
		jen.Id("value").Interface(),
	)
}

func (g *QuerysetGenerator) AddConstructor() {
	g.File.Func().Id(g.names().QuerysetStructConstructor).Params(
	// No args
	).Params(jen.Op("*").Id(g.names().QuerysetStruct)).Block(
		jen.Return(
			jen.Op("&").Id(g.names().QuerysetStruct).Values(jen.Dict{
				jen.Id("filter"):  jen.Id("make").Call(jen.Index().Id(g.names().QuerysetFilterArgStruct), jen.Lit(0)),
				jen.Id("orderBy"): jen.Id("make").Call(jen.Index().Id(g.names().QuerysetOrderByArgStruct), jen.Lit(0)),
				jen.Id("limit"):   jen.Lit(0),
				jen.Id("offset"):  jen.Lit(0),
			}),
		),
	)
}

func (g *QuerysetGenerator) AddFilterMethod() {
	g.AddChainedMethod("Filter",
		[]jen.Code{jen.Id("filter").Op("...").Id(g.names().QuerysetFilterArgStruct)},
		[]jen.Code{
			jen.Id("qs").Dot("filter").Op("=").Id("append").Params(
				jen.Id("qs").Dot("filter"),
				jen.Id("filter").Op("..."),
			),
			jen.Return(jen.Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddOrderMethod() {
	g.AddChainedMethod("Order",
		[]jen.Code{jen.Id("orderBy").Op("...").Id(g.names().QuerysetOrderByArgStruct)},
		[]jen.Code{
			jen.Id("qs").Dot("orderBy").Op("=").Id("append").Params(
				jen.Id("qs").Dot("orderBy"),
				jen.Id("orderBy").Op("..."),
			),
			jen.Return(jen.Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddLimitMethod() {
	g.AddChainedMethod("Limit",
		[]jen.Code{jen.Id("limit").Uint64()},
		[]jen.Code{
			jen.Id("qs").Dot("limit").Op("=").Id("limit"),
			jen.Return(jen.Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddOffsetMethod() {
	g.AddChainedMethod("Offset",
		[]jen.Code{jen.Id("offset").Uint64()},
		[]jen.Code{
			jen.Id("qs").Dot("offset").Op("=").Id("offset"),
			jen.Return(jen.Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddStruct() {
	g.File.Type().Id(g.names().QuerysetStruct).Struct(
		jen.Id("filter").Index().Id(g.names().QuerysetFilterArgStruct),
		jen.Id("orderBy").Index().Id(g.names().QuerysetOrderByArgStruct),
		jen.Id("limit").Uint64(),
		jen.Id("offset").Uint64(),
	)
}

// AddChainedMethod is a helper to add a struct method that returns an instance
// of the struct for chaining.
func (g *QuerysetGenerator) AddChainedMethod(name string, args []jen.Code, block []jen.Code) {
	g.File.Func().Params(
		jen.Id("qs").Op("*").Id(g.names().QuerysetStruct),
	).Id(name).Params(
		args...,
	).Params(jen.Op("*").Id(g.names().QuerysetStruct)).Block(
		block...,
	)
}

func (g *QuerysetGenerator) AddFilterOptionsStruct() {
	fields := make([]jen.Code, 0)
	values := jen.Dict{}
	for _, f := range g.Model.Fields() {
		fieldName := f.Settings().Name
		fieldType := f.Settings().Names(g.Model.Settings().Name).FilterOptionStruct
		fields = append(fields, jen.Id(fieldName).Id(fieldType))
		values[jen.Id(fieldName)] = jen.Op("&").Id(fieldType).Values()
	}

	// Define struct
	g.File.Type().Id(g.names().QuerysetFilterOptionsStruct).Struct(fields...)

	// Create instance of struct
	varName := g.names().QuerysetFilterOptionsStruct
	g.File.Var().Id(g.names().FilterOptionsVar).Op("=").
		Op("&").Id(varName).Values(values)
}

func (g *QuerysetGenerator) Generate() {
	// Create main struct and constructor
	g.AddStruct()
	g.AddConstructor()

	// Methods
	g.AddFilterMethod()
	g.AddOrderMethod()
	g.AddLimitMethod()
	g.AddOffsetMethod()
	// TODO: g.AddUpdateMethod()
	// TODO: g.AddAllMethod()
	// TODO: g.AddOneMethod()
	// TODO: g.AddDeleteMethod()
	// TODO: g.AddCountMethod()
	// TODO: g.AddScanMethod()
	// TODO: g.AddStarSelectMethod()

	// Other types
	g.AddFilterArgsStruct()
	g.AddOrderByArgsStruct()
	g.AddSetterArgsStruct()

	// Where helper
	g.AddFilterOptionsStruct()
}
