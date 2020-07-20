package banister

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
)

type ManagerGenerator struct {
	Model Model
	File  *File
}

func NewManagerGenerator(file *File, model Model) *ManagerGenerator {
	return &ManagerGenerator{Model: model, File: file}
}

func (g *ManagerGenerator) names() GeneratedModelNames {
	return g.Model.Settings().Names()
}

// AddMethod is a helper to add a struct method
func (g *ManagerGenerator) AddMethod(name string, args []Code, block []Code, returns []Code) {
	receiver := Id("mgr").Op("*").Id(g.names().ManagerStruct)
	g.File.Func().Params(receiver).Id(name).Params(args...).Params(returns...).Block(block...)
}

func (g *ManagerGenerator) AddStruct() {
	g.File.Type().Id(g.names().ManagerStruct).Struct(
		Id("db").Op("*").Qual("database/sql", "DB"),
		Id("storeConfig").Id(globalNames.StoreConfigStruct),
		Id("config").Id(g.names().ConfigStruct),
	)
}

func (g *ManagerGenerator) AddConstructor() {
	g.File.Func().Id(g.names().ManagerConstructor).Params(
		Id("db").Op("*").Qual("database/sql", "DB"),
		Id("storeConfig").Id(globalNames.StoreConfigStruct),
		Id("config").Id(g.names().ConfigStruct),
	).Params(Op("*").Id(g.names().ManagerStruct)).Block(
		Return(
			Op("&").Id(g.names().ManagerStruct).Values(Dict{
				Id("db"):          Id("db"),
				Id("storeConfig"): Id("storeConfig"),
				Id("config"):      Id("config"),
			}),
		),
	)
}

func (g *ManagerGenerator) AddDeleteMethod() {
	g.AddMethod("Delete",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{
			Panic(Lit("implement me")),
		},
		[]Code{Error()},
	)
}

func (g *ManagerGenerator) AddAllMethod() {
	g.AddMethod("All",
		[]Code{},
		[]Code{Return(Op("mgr").Dot("Filter").Call().Dot("All").Call())},
		[]Code{Index().Id(g.names().ModelStruct), Error()},
	)
}

func (g *ManagerGenerator) AddInsertMethod() {
	g.AddMethod("Insert",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Error()},
	)
}

func (g *ManagerGenerator) AddInsertInstanceMethod() {
	g.AddMethod("insertInstance",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{
			Panic(Lit("implement me")),
		},
		[]Code{Error()},
	)
}

func (g *ManagerGenerator) AddUpdateMethod() {
	g.AddMethod("Update",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Error()},
	)
}

func (g *ManagerGenerator) AddGetMethod() {
	pkGoType := fmt.Sprintf("%T", PrimaryKeyField(g.Model).EmptyDefault())
	g.AddMethod("Get",
		[]Code{Id("id").Id(pkGoType)},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Op("*").Id(g.names().ModelStruct)},
	)
}

func (g *ManagerGenerator) AddOrMethod() {
	g.AddMethod("Or",
		[]Code{Id("filter").Op("...").Id(g.names().QuerysetFilterArgStruct)},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Op("*").Id(g.names().QuerysetStruct)},
	)
}

func (g *ManagerGenerator) AddAndMethod() {
	g.AddMethod("And",
		[]Code{Id("filter").Op("...").Id(g.names().QuerysetFilterArgStruct)},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Op("*").Id(g.names().QuerysetStruct)},
	)
}

func (g *ManagerGenerator) AddNewModelMethod() {
	g.AddMethod("newModel",
		[]Code{},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Op("*").Id(g.names().ModelStruct)},
	)
}

func (g *ManagerGenerator) AddFilterMethod() {
	g.AddMethod(
		"Filter",
		[]Code{Id("filter").Op("...").Id(g.names().QuerysetFilterArgStruct)},
		[]Code{
			Id("v").Op(":=").Id(g.names().QuerysetConstructor).Call(Id("mgr")),
			Id("v").Dot("Filter").Call(Id("filter").Op("...")),
			Return(Id("v")),
		},
		[]Code{Op("*").Id(g.names().QuerysetStruct)},
	)
}

func (g *ManagerGenerator) Generate() {
	g.AddStruct()
	g.AddConstructor()
	g.AddFilterMethod()
	g.AddDeleteMethod()
	g.AddInsertMethod()
	g.AddInsertInstanceMethod()
	g.AddUpdateMethod()
	g.AddAllMethod()
	g.AddGetMethod()
	g.AddNewModelMethod()
	g.AddAndMethod()
	g.AddOrMethod()
}
