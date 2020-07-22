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
func (g *ManagerGenerator) AddMethod(name string, args, block, returns []Code) {
	receiver := Id("mgr").Op("*").Id(g.names().ManagerStruct)
	g.File.Func().Params(receiver).Id(name).
		Params(args...).
		Params(returns...).
		Block(block...)
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
	pkSettings := PrimaryKeyField(g.Model).Settings()
	callQuerysetDelete := Err().Op(":=").Id("mgr").Dot("Filter").Call(
		Id(g.names().QuerysetFilterArgStruct).Values(Dict{
			Id("filter"): Op("&").Qual("github.com/Masterminds/squirrel", "Eq").Values(Dict{
				Lit(pkSettings.DBColumn): Id("m").Dot(pkSettings.Name),
			}),
		}),
	).Dot("Delete").Call()

	checkErr := If(Err().Op("!=").Nil()).Block(Return(Err()))

	g.AddMethod("Delete",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{
			g.MaybeCallHook(globalNames.HookPreDelete).Line(),
			Comment("Call delete on queryset with PK as the filter"),
			callQuerysetDelete,
			checkErr.Line(),
			g.MaybeCallHook(globalNames.HookPostDelete).Line(),
			Return(Nil()),
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

func (g *ManagerGenerator) MaybeCallHook(hookName string) *Statement {
	return Comment("Call hook if provided").Line().If(
		Id("mgr").Dot("config").Dot(hookName).Op("!=").Nil(),
	).Block(
		Id("mgr").Dot("config").Dot(hookName).Call(Id("m")),
	)
}

func (g *ManagerGenerator) AddInsertInstanceMethod() {
	columns := make([]Code, 0)
	values := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		// Skip auto-created fields
		if f.Type() == Auto {
			continue
		}

		columns = append(columns, Lit(f.Settings().DBColumn))
		values = append(values, Op("m").Dot(f.Settings().Name))
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

	g.AddMethod("insertInstance",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{
			g.MaybeCallHook(globalNames.HookPreInsert).Line(),
			defineQuery,
			addColumns,
			addValues.Line(),
			toSQL.Line(),
			exec,
			checkErr.Line(),
			lastInsertId,
			checkErr.Line(),
			assignPK.Line(),
			g.MaybeCallHook(globalNames.HookPostInsert).Line(),
			Return(Nil()),
		},
		[]Code{Error()},
	)
}

func (g *ManagerGenerator) AddUpdateMethod() {
	checkErr := If(Err().Op("!=").Nil()).Block(Return(Err()))
	setters := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		// TODO: Maybe don't update primary key? Should see what Django
		//   does here.
		setter := Id(g.names().QuerysetSetterArgStruct).Values(Dict{
			Id("field"): Lit(f.Settings().DBColumn),
			Id("value"): Id("m").Dot(f.Settings().Name),
		})
		setters = append(setters, setter)
	}

	pkSettings := PrimaryKeyField(g.Model).Settings()
	callQuerysetDelete := Err().Op(":=").Id("mgr").Dot("Filter").Call(
		Id(g.names().QuerysetFilterArgStruct).Values(Dict{
			Id("filter"): Op("&").Qual("github.com/Masterminds/squirrel", "Eq").Values(Dict{
				Lit(pkSettings.DBColumn): Id("m").Dot(pkSettings.Name),
			}),
		}),
	).Dot("Update").Call(setters...)

	g.AddMethod("Update",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{
			g.MaybeCallHook(globalNames.HookPreUpdate).Line(),
			Comment("Call update on queryset with PK as the filter"),
			callQuerysetDelete,
			checkErr.Line(),
			g.MaybeCallHook(globalNames.HookPostUpdate).Line(),
			Return(Nil()),
		},
		[]Code{Error()},
	)
}

func (g *ManagerGenerator) AddGetMethod() {
	pkGoType := fmt.Sprintf("%T", PrimaryKeyField(g.Model).EmptyDefault())
	g.AddMethod("Get",
		[]Code{Id("id").Id(pkGoType)},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Op("*").Id(g.names().ModelStruct), Error()},
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

		goType := fmt.Sprintf("%T", f.EmptyDefault())
		goDefaultVal := f.Settings().Default.Value

		// If no default is provided and nil is not allowed, use the
		// fallback default
		if f.Settings().Null == false && goDefaultVal == nil {
			goDefaultVal = f.EmptyDefault()
		}

		switch v := goDefaultVal.(type) {
		case nil:
			defaultValue = Nil()
			castTo = Id(goType)
		case time.Time:
			defaultValue = Qual("time", "Time").Values()
			castTo = Qual("time", "Time")
		case time.Duration:
			defaultValue = Qual("time", "Duration").Call()
			castTo = Qual("time", "Duration")
		default:
			defaultValue = Lit(v)
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
					Id("s").Dot("value").Op("!=").Nil(),
				).Block(
					Id("m").Dot(f.Settings().Name).Op("=").Id("s").Dot("value").Op(".").Params(castTo),
				).Else().Block(
					Id("m").Dot(f.Settings().Name).Op("=").Nil().Comment("Cannot cast nil"),
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

	// Add default case that panics
	setterCases = append(setterCases, Default().Block(
		Panic(Lit("invalid field for setter: ").Op("+").Id("s").Dot("field")),
	))

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
