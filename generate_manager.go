package banister

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"time"
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
	createInstance := Id("m").Op(":=").Id("mgr").Dot("newModel").Call(Id("set").Op("..."))
	insert := Err().Op(":=").Id("mgr").Dot("insertInstance").Call(Id("m"))
	checkErr := If(Err().Op("!=").Nil()).Block(Return(Nil(), Err()))
	g.AddMethod("Insert",
		[]Code{Id("set").Op("...").Id(g.names().QuerysetSetterArgStruct)},
		[]Code{
			createInstance,
			insert,
			checkErr,
			Return(Id("m"), Err()),
		},
		[]Code{Op("*").Id(g.names().ModelStruct), Error()},
	)
}

func (g *ManagerGenerator) AddInsertInstanceMethod() {
	callPreInsertHook := If(
		Id("mgr").Dot("config").Dot(globalNames.HookPreInsert).Op("!=").Nil(),
	).Block(
		Id("mgr").Dot("config").Dot(globalNames.HookPreInsert).Call(Id("m")),
	)

	columns := make([]Code, 0)
	values := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		columns = append(columns, Line().Lit(f.Settings().DBColumn))
		values = append(values, Line().Op("m").Dot(f.Settings().Name))
	}

	defineQuery := Id("query").Op(":=").Qual("github.com/Masterminds/squirrel", "Insert").Call(
		Lit(g.Model.Settings().DBTable),
	)

	addColumns := Id("query").Op("=").Id("query").Dot("Columns").Call(columns...)
	addValues := Id("query").Op("=").Id("query").Dot("Values").Call(values...)
	toSQL := Id("q").Op(",").Id("args").Op(":=").Id("mgr").Dot("toSQL").Call(Id("query"))

	exec := Id("result").Op(",").Err().Op(":=").
		Id("mgr").Dot("db").Dot("Exec").Call(Id("q"), Id("args").Op("..."))

	checkErr := If(Err().Op("!=").Nil()).Block(Return(Err()))

	lastInsertId := Id("id").Op(",").Err().Op(":=").Id("result").Dot("LastInsertId").Call()

	assignPK := Comment("Update PK on model").Line().
		Id("m").Dot(PrimaryKeyField(g.Model).Settings().Name).Op("=").Id("id")

	callPostInsertHook := If(
		Id("mgr").Dot("config").Dot(globalNames.HookPostInsert).Op("!=").Nil(),
	).Block(
		Id("mgr").Dot("config").Dot(globalNames.HookPostInsert).Call(Id("m")),
	)

	g.AddMethod("insertInstance",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{
			callPreInsertHook.Line(),
			defineQuery.Line(),
			addColumns.Line(),
			addValues.Line(),
			toSQL.Line(),
			exec,
			checkErr.Line(),
			lastInsertId,
			checkErr.Line(),
			assignPK.Line(),
			callPostInsertHook.Line(),
			Return(Nil()),
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

func (g *ManagerGenerator) AddToSQLMethod() {
	genSQL := Id("query").Op(",").Id("args").Op(",").Err().Op(":=").
		Id("q").Dot("ToSql").Call()

	maybePanic := If(
		Err().Op("!=").Nil(),
	).Block(
		Panic(Err()),
	)

	g.AddMethod("toSQL",
		[]Code{Id("q").Qual("github.com/Masterminds/squirrel", "Sqlizer")},
		[]Code{
			genSQL,
			maybePanic,
			Return(Id("query"), Id("args")),
		},
		[]Code{
			String(),
			Index().Interface(),
		},
	)
}

func (g *ManagerGenerator) AddNewModelMethod() {
	defaultValues := Dict{}
	setterCases := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		var castTo Code
		var defaultValue Code

		switch v := f.EmptyDefault().(type) {
		case time.Time:
			defaultValue = Qual("time", "Time").Values()
			castTo = Qual("time", "Time")
		case time.Duration:
			defaultValue = Qual("time", "Duration").Call()
			castTo = Qual("time", "Duration")
		default:
			defaultValue = Lit(v)
			goType := fmt.Sprintf("%T", f.EmptyDefault())
			castTo = Id(goType)
		}

		// Nullable fields have null-friendly defaults and pointer casts
		if f.Settings().Null {
			castTo = Id("*").Add(castTo)
			defaultValue = Nil()
		}

		// Add setter to cases
		var newCase Code
		if f.Settings().Null {
			newCase = Case(Lit(f.Settings().DBColumn)).Block(
				If(
					Id("s").Dot("value").Op("==").Nil(),
				).Block(
					Id("m").Dot(f.Settings().Name).Op("=").Nil(),
				).Else().Block(
					Id("m").Dot(f.Settings().Name).Op("=").
						Id("s").Dot("value").Op(".").Params(castTo),
				),
			)
		} else {
			newCase = Case(Lit(f.Settings().DBColumn)).Block(
				Id("m").Dot(f.Settings().Name).Op("=").
					Id("s").Dot("value").Op(".").Params(castTo),
			)
		}

		defaultValues[Id(f.Settings().Name)] = defaultValue
		setterCases = append(setterCases, newCase)
	}

	instantiateWithDefaults := Id("m").Op(":=").
		Id(g.names().ModelStruct).Values(defaultValues)

	loopOverSetters := For(
		Op("_").Op(",").Id("s").Op(":=").Range().Id("set"),
	).Block(
		Switch(Id("s").Dot("field")).Block(setterCases...),
	)

	g.AddMethod("newModel",
		[]Code{Id("set").Op("...").Id(g.names().QuerysetSetterArgStruct)},
		[]Code{
			instantiateWithDefaults.Line(),
			loopOverSetters.Line(),
			Return(Op("&").Id("m")),
		},
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
	g.AddToSQLMethod()
}
