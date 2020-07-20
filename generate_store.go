package banister

import (
	"github.com/dave/jennifer/jen"
)

type StoreGenerator struct {
	Models []Model
	File   *jen.File
}

func NewStoreGenerator(file *jen.File, models []Model) *StoreGenerator {
	return &StoreGenerator{Models: models, File: file}
}

func (g *StoreGenerator) AddStruct() {
	g.File.Type().Id(globalNames.StoreStruct).Struct(
		jen.Id("db").Op("*").Qual("database/sql", "Conn"),
		jen.Id("config").Id(globalNames.StoreConfigStruct),
	)
}

func (g *StoreGenerator) AddConfigStruct() {
	modelConfigs := make([]jen.Code, 0)
	for _, m := range g.Models {
		name := PrivateGoName(m.Settings().Names().ConfigStruct)
		modelConfigs = append(modelConfigs, jen.Id(name).Id(m.Settings().Names().ConfigStruct))
	}

	g.File.Type().Id(globalNames.StoreConfigStruct).Struct(
		append([]jen.Code{
			jen.Id("connectionStr").String(),
			jen.Id("maxIdleConnections").Int(),
			jen.Id("maxOpenConnections").Int(),
			jen.Id("connectionMaxLifetime").Qual("time", "Duration"),
			jen.Line().Comment("Model configs").Line(),
		}, modelConfigs...)...,
	)
}

func (g *StoreGenerator) Generate() {
	g.AddStruct()
	g.AddConfigStruct()
}
