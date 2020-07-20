package banister

import (
	. "github.com/dave/jennifer/jen"
)

type StoreGenerator struct {
	Models []Model
	File   *File
}

func NewStoreGenerator(file *File, models []Model) *StoreGenerator {
	return &StoreGenerator{Models: models, File: file}
}

func (g *StoreGenerator) AddStruct() {
	g.File.Type().Id(globalNames.StoreStruct).Struct(
		Id("db").Op("*").Qual("database/sql", "Conn"),
		Id("config").Id(globalNames.StoreConfigStruct),
	)
}

func (g *StoreGenerator) AddConfigStruct() {
	modelConfigs := make([]Code, 0)
	for _, m := range g.Models {
		name := PrivateGoName(m.Settings().Names().ConfigStruct)
		modelConfigs = append(modelConfigs, Id(name).Id(m.Settings().Names().ConfigStruct))
	}

	g.File.Type().Id(globalNames.StoreConfigStruct).Struct(
		append([]Code{
			Id("connectionStr").String(),
			Id("maxIdleConnections").Int(),
			Id("maxOpenConnections").Int(),
			Id("connectionMaxLifetime").Qual("time", "Duration"),
			Line().Comment("Model configs").Line(),
		}, modelConfigs...)...,
	)
}

func (g *StoreGenerator) Generate() {
	g.AddStruct()
	g.AddConfigStruct()
}
