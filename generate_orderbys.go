package banister

import (
	. "github.com/dave/jennifer/jen"
)

type OrderBysGenerator struct {
	Model Model
	File  *File
}

func NewOrderBysGenerator(file *File, model Model) *OrderBysGenerator {
	return &OrderBysGenerator{Model: model, File: file}
}

func (g *OrderBysGenerator) AddOrderStructs() {
	for _, f := range g.Model.Fields() {
		structName := f.Settings().Names(g.Model).QuerysetOrderByDirectionStruct
		g.File.Type().Id(structName).Struct(
			Id("Asc").Id(g.Model.Settings().Names().QuerysetOrderByArgStruct),
			Id("Desc").Id(g.Model.Settings().Names().QuerysetOrderByArgStruct),
		)
	}
}

func (g *OrderBysGenerator) AddStruct() {
	fields := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		name := f.Settings().Name
		structName := f.Settings().Names(g.Model).QuerysetOrderByDirectionStruct
		fields = append(fields, Id(name).Id(structName))
	}

	g.File.Type().Id(g.Model.Settings().Names().QuerysetOrderByOptionsStruct).Struct(fields...)
}

func (g *OrderBysGenerator) AddConstructor() {
	values := Dict{}
	for _, f := range g.Model.Fields() {
		asc := Dict{
			Id("field"): Lit(f.Settings().Names(g.Model).QualifiedColumn),
			Id("order"): Lit("ASC"),
			Id("join"):  Lit(""),
		}
		desc := Dict{
			Id("field"): Lit(f.Settings().Names(g.Model).QualifiedColumn),
			Id("order"): Lit("ASC"),
			Id("join"):  Lit(""),
		}
		values[Id(f.Settings().Name)] = Id(f.Settings().Names(g.Model).QuerysetOrderByDirectionStruct).Values(Dict{
			Id("Asc"):  Id(g.Model.Settings().Names().QuerysetOrderByArgStruct).Values(asc),
			Id("Desc"): Id(g.Model.Settings().Names().QuerysetOrderByArgStruct).Values(desc),
		})
	}

	names := g.Model.Settings().Names()
	g.File.Func().Id(names.QuerysetOrderByOptionsConstructor).Params(
	// No args
	).Params(
		Id(names.QuerysetOrderByOptionsStruct),
	).Block(
		Return(Id(names.QuerysetOrderByOptionsStruct).Values(values)),
	)
}

func (g *OrderBysGenerator) Generate() {
	g.AddStruct()
	g.AddOrderStructs()
	g.AddConstructor()

}
