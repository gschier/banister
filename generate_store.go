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
	modelConfigs := make([]Code, 0)
	for _, m := range g.Models {
		name := m.Settings().Names().ManagerAccessor
		fieldDef := Id(name).Op("*").Id(m.Settings().Names().ManagerStruct)
		modelConfigs = append(modelConfigs, fieldDef)
	}

	g.File.Type().Id(globalNames.StoreStruct).Struct(
		append([]Code{
			Id("db").Op("*").Qual("database/sql", "DB"),
			Id("config").Id(globalNames.StoreConfigStruct),
			Line().Comment("Managers").Line(),
		}, modelConfigs...)...,
	)
}

func (g *StoreGenerator) AddConstructor() {
	values := Dict{
		Id("db"):     Id("db"),
		Id("config"): Id("c"),
	}

	for _, m := range g.Models {
		names := m.Settings().Names()
		values[Id(names.ManagerAccessor)] = Id(names.ManagerConstructor).Call(
			Id("db"),
			Id("c"),
			Id("c").Dot(names.ConfigStruct),
		)
	}

	g.File.Func().Id(globalNames.StoreConstructor).Params(
		Id("db").Op("*").Qual("database/sql", "DB"),
		Id("c").Id(globalNames.StoreConfigStruct),
	).Params(Op("*").Id(globalNames.StoreStruct)).Block(
		Return(Op("&").Id(globalNames.StoreStruct).Values(values)),
	)
}

func (g *StoreGenerator) AddConfigStruct() {
	modelConfigs := make([]Code, 0)
	for _, m := range g.Models {
		name := m.Settings().Names().ConfigStruct
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
	g.AddConstructor()
	g.AddConfigStruct()
}
