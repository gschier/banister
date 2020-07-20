package banister

import (
	"github.com/dave/jennifer/jen"
)

type ManagerGenerator struct {
	Model Model
	File  *jen.File
}

func NewManagerGenerator(file *jen.File, model Model) *ManagerGenerator {
	return &ManagerGenerator{Model: model, File: file}
}

func (g *ManagerGenerator) names() GeneratedModelNames {
	return g.Model.Settings().Names()
}

// AddMethod is a helper to add a struct method
func (g *ManagerGenerator) AddMethod(name string, args []jen.Code, block []jen.Code, returns []jen.Code) {
	receiver := jen.Id("mgr").Op("*").Id(g.names().ManagerStruct)
	g.File.Func().Params(receiver).Id(name).Params(args...).Params(returns...).Block(block...)
}

func (g *ManagerGenerator) AddStruct() {
	g.File.Type().Id(g.names().ManagerStruct).Struct(
		jen.Id("db").Qual("database/sql", "DB"),
		jen.Id("storeConfig").Id(globalNames.StoreConfigStruct),
		jen.Id("config").Id(g.names().ConfigStruct),
	)
}

func (g *ManagerGenerator) AddConstructor() {
	g.File.Func().Id(g.names().ManagerConstructor).Params(
		jen.Id("db").Op("*").Qual("database/sql", "DB"),
		jen.Id("storeConfig").Id(globalNames.StoreConfigStruct),
		jen.Id("config").Id(g.names().ConfigStruct),
	).Params(jen.Op("*").Id(g.names().ManagerStruct)).Block(
		jen.Return(
			jen.Op("&").Id(g.names().ManagerStruct).Values(jen.Dict{
				jen.Id("db"):          jen.Id("db"),
				jen.Id("storeConfig"): jen.Id("storeConfig"),
				jen.Id("config"):      jen.Id("config"),
			}),
		),
	)
}

func (g *ManagerGenerator) AddFilterMethod() {
	g.AddMethod(
		"Filter",
		[]jen.Code{
			jen.Id("filter").Op("...").Id(g.names().QuerysetFilterArgStruct),
		},
		[]jen.Code{
			jen.Id("v").Op(":=").Id(g.names().QuerysetConstructor).Call(jen.Id("mgr")),
			jen.Id("v").Dot("Filter").Call(jen.Id("filter").Op("...")),
			jen.Return(jen.Id("v")),
		},
		[]jen.Code{
			jen.Op("*").Id(g.names().QuerysetStruct),
		},
	)
}

func (g *ManagerGenerator) Generate() {
	g.AddStruct()
	g.AddConstructor()
	g.AddFilterMethod()
}
