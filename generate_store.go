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

	g.File.Comment(globalNames.StoreStruct + " defines a data store")
	g.File.Type().Id(globalNames.StoreStruct).Struct(
		append([]Code{
			Id("DB").Op("*").Qual("database/sql", "DB"),
			Id("config").Id(globalNames.StoreConfigStruct),
			Line().Comment("model managers").Line(),
		}, modelConfigs...)...,
	)
}

func (g *StoreGenerator) AddConstructor() {
	values := Dict{
		Id("DB"):     Id("db"),
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

	g.File.Comment(globalNames.StoreConstructor + " returns a new store instance")
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
		configDef := Id(name).Id(m.Settings().Names().ConfigStruct)
		modelConfigs = append(modelConfigs, configDef)
	}

	g.File.Comment(globalNames.StoreConfigStruct + " holds configuration for the store")
	g.File.Type().Id(globalNames.StoreConfigStruct).Struct(
		append([]Code{
			Id("connectionStr").String(),
			Id("maxIdleConnections").Int(),
			Id("maxOpenConnections").Int(),
			Id("connectionMaxLifetime").Qual("time", "Duration"),
			Line().Comment("model configs").Line(),
		}, modelConfigs...)...,
	)
}

func (g *StoreGenerator) AddGlobalWhere() {
	fields := make([]Code, 0)
	for _, m := range g.Models {
		structName := m.Settings().Names().QuerysetFilterOptionsStruct
		fields = append(fields, Id(m.Settings().Name).Id(structName))
	}
	g.File.Comment("Where contains helpers for filtering querysets")
	g.File.Var().Id("Where").Op("=").Struct(fields...).Values()
}

func (g *StoreGenerator) AddGlobalSetters() {
	fields := make([]Code, 0)
	for _, m := range g.Models {
		structName := m.Settings().Names().QuerysetSetterOptionsStruct
		fields = append(fields, Id(m.Settings().Name).Id(structName))
	}

	g.File.Comment("Set contains helpers for setting fields during inserts and updates")
	g.File.Var().Id("Set").Op("=").Struct(fields...).Values()
}

func (g *StoreGenerator) AddGlobalOrderBys() {
	fields := make([]Code, 0)
	values := Dict{}

	for _, m := range g.Models {
		name := m.Settings().Name
		names := m.Settings().Names()
		fields = append(fields, Id(name).Id(names.QuerysetOrderByOptionsStruct))
		values[Id(name)] = Id(names.QuerysetOrderByOptionsConstructor).Call()
	}

	g.File.Comment("OrderBy contains helpers for sorting query results")
	g.File.Var().Id("OrderBy").Op("=").Struct(fields...).Values(values)
}

func (g *StoreGenerator) Generate() {
	g.AddStruct()
	g.AddConstructor()
	g.AddConfigStruct()

	g.AddGlobalWhere()
	g.AddGlobalSetters()
	g.AddGlobalOrderBys()
}
