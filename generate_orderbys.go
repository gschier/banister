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

func (g *OrderBysGenerator) Generate() {
	fields := make([]Code, 0)
	values := Dict{}
	for _, f := range g.Model.Fields() {
		name := f.Settings().Name + "Ascending"
		structName := g.Model.Settings().Names().QuerysetOrderByArgStruct
		fields = append(fields, Id(name).Id(structName))
		values[Id(name)] = Id(structName).Values(Dict{
			Id("field"): Lit(f.Settings().Names(g.Model).QualifiedColumn),
			Id("order"): Lit("ASC"),
			Id("join"):  Lit(""),
		})

		name = f.Settings().Name + "Descending"
		structName = g.Model.Settings().Names().QuerysetOrderByArgStruct
		fields = append(fields, Id(name).Id(structName))
		values[Id(name)] = Id(structName).Values(Dict{
			Id("field"): Lit(f.Settings().Names(g.Model).QualifiedColumn),
			Id("order"): Lit("DESC"),
			Id("join"):  Lit(""),
		})
	}

	g.File.Var().Id(g.Model.Settings().Names().OrderByOptionsVar).Op("=").
		Struct(fields...).Values(values)
}
