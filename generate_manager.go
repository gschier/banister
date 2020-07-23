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

func (g *ManagerGenerator) AddMethodWithPanicVariant(name string, args, panicArgs, block, returns []Code) {
	g.AddMethod(name, args, block, returns)

	// Now we're going to add a variant of the method that panics instead of
	// returning an error

	var (
		panicVariantName          = name + "P"
		callOriginalAndMaybePanic []Code
		panicVariantReturns       []Code
	)

	if len(returns) == 1 {
		// Original method only returned error, so only expect error back
		panicVariantReturns = []Code{}
		callOriginalAndMaybePanic = []Code{
			Err().Op(":=").Id("mgr").Dot(name).Call(panicArgs...),
			If(Err().Op("!=").Nil()).Block(Panic(Err())),
		}
	} else {
		// Original method returned a value too, so handle that as well
		panicVariantReturns = []Code{returns[0]}
		callOriginalAndMaybePanic = []Code{
			Id("v").Op(",").Err().Op(":=").Id("mgr").Dot(name).Call(panicArgs...),
			If(Err().Op("!=").Nil()).Block(Panic(Err())),
			Return(Id("v")),
		}
	}

	g.AddMethod(panicVariantName, args, callOriginalAndMaybePanic, panicVariantReturns)
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
	checkErr := If(Err().Op("!=").Nil()).Block(Return(Err()))

	filterDef := Op("&").Qual("github.com/Masterminds/squirrel", "Eq").Values(Dict{
		Lit(pkSettings.DBColumn): Id("m").Dot(pkSettings.Name),
	})

	callQuerysetDelete := Err().Op(":=").Id("mgr").Dot("Filter").Call(
		Id(g.names().QuerysetFilterArgStruct).Values(Dict{Id("filter"): filterDef}),
	).Dot("Delete").Call()

	g.AddMethodWithPanicVariant("Delete",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{Id("m")},
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
	g.AddMethodWithPanicVariant("All",
		[]Code{},
		[]Code{},
		[]Code{Return(Op("mgr").Dot("Filter").Call().Dot("All").Call())},
		[]Code{Index().Id(g.names().ModelStruct), Error()},
	)
}

func (g *ManagerGenerator) AddInsertMethod() {
	createInstance := Id("m").Op(":=").Id("mgr").Dot("newModel").Call(Id("set").Op("..."))
	insert := Err().Op(":=").Id("mgr").Dot("insertInstance").Call(Id("m"))
	checkErr := If(Err().Op("!=").Nil()).Block(Return(Nil(), Err()))
	g.AddMethodWithPanicVariant("Insert",
		[]Code{Id("set").Op("...").Id(g.names().QuerysetSetterArgStruct)},
		[]Code{Id("set").Op("...")},
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
		if f.Type() == Auto {
			continue // Skip auto-created fields
		}
		columns = append(columns, Lit(f.Settings().DBColumn))
		values = append(values, Op("m").Dot(f.Settings().Name))
	}

	insert := Id("query").Op(":=").Qual("github.com/Masterminds/squirrel", "Insert")
	defineQuery := insert.Call(Lit(g.Model.Settings().DBTable))

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
	filterDef := Op("&").Qual("github.com/Masterminds/squirrel", "Eq").Values(Dict{
		Lit(pkSettings.DBColumn): Id("m").Dot(pkSettings.Name),
	})

	callQuerysetDelete := Err().Op(":=").Id("mgr").Dot("Filter").Call(
		Id(g.names().QuerysetFilterArgStruct).Values(Dict{Id("filter"): filterDef}),
	).Dot("Update").Call(setters...)

	g.AddMethodWithPanicVariant("Update",
		[]Code{Id("m").Op("*").Id(g.names().ModelStruct)},
		[]Code{Id("m")},
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
	pk := PrimaryKeyField(g.Model)
	pkGoType := fmt.Sprintf("%T", pk.EmptyDefault())

	eqDef := Qual("github.com/Masterminds/squirrel", "Eq").Values(Dict{
		Lit(pk.Settings().Names(g.Model).QualifiedColumn): Id("id"),
	})

	callFilterAndOne := Id("mgr").Dot("Filter").Call(
		Id(g.names().QuerysetFilterArgStruct).Values(Dict{Id("filter"): eqDef}),
	).Dot("One").Call()

	g.AddMethodWithPanicVariant("Get",
		[]Code{Id("id").Id(pkGoType)},
		[]Code{Id("id")},
		[]Code{Return(callFilterAndOne)},
		[]Code{Op("*").Id(g.names().ModelStruct), Error()},
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
	setterCases := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		var castTo *Statement

		goType := fmt.Sprintf("%T", f.EmptyDefault())
		goDefaultVal := f.Settings().Default.Value

		// If no default is provided and nil is not allowed, use the
		// fallback default
		if f.Settings().Null == false && goDefaultVal == nil {
			goDefaultVal = f.EmptyDefault()
		}

		switch goDefaultVal.(type) {
		case nil:
			castTo = Id(goType)
		case time.Time:
			castTo = Qual("time", "Time")
		case time.Duration:
			castTo = Qual("time", "Duration")
		default:
			castTo = Id(goType)
		}

		// Nullable fields have null-friendly defaults and pointer casts
		if f.Settings().Null {
			castTo = Id("*").Add(castTo)
		}

		// Add setter to cases
		var newCase Code
		if f.Settings().Null {
			newCase = Case(Lit(f.Settings().DBColumn)).Block(
				If(Id("s").Dot("value").Op("!=").Nil()).Block(
					Id("m").Dot(f.Settings().Name).Op("=").Id("s").Dot("value").Op(".").Params(castTo),
				).Else().Block(
					Id("m").Dot(f.Settings().Name).Op("=").Nil().Comment("Cannot cast nil"),
				),
			)
		} else {
			newCase = Case(Lit(f.Settings().DBColumn)).Block(
				Id("m").Dot(f.Settings().Name).Op("=").Id("s").Dot("value").Op(".").Params(castTo),
			)
		}

		setterCases = append(setterCases, newCase)
	}

	// Add default case that panics
	setterCases = append(setterCases, Default().Block(
		Panic(Lit("invalid field for setter: ").Op("+").Id("s").Dot("field")),
	))

	instantiate := Id("m").Op(":=").Id(g.names().ModelConstructor).Call()

	loopOverSetters := For(
		Op("_").Op(",").Id("s").Op(":=").Range().Id("set"),
	).Block(
		Switch(Id("s").Dot("field")).Block(setterCases...),
	)

	g.AddMethod("newModel",
		[]Code{Id("set").Op("...").Id(g.names().QuerysetSetterArgStruct)},
		[]Code{
			instantiate.Line(),
			Comment("Apply setters to default struct"),
			loopOverSetters.Line(),
			Return(Id("m")),
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

func (g *ManagerGenerator) AddExcludeMethod() {
	g.AddMethod(
		"Exclude",
		[]Code{Id("exclude").Op("...").Id(g.names().QuerysetFilterArgStruct)},
		[]Code{
			Id("v").Op(":=").Id(g.names().QuerysetConstructor).Call(Id("mgr")),
			Id("v").Dot("Exclude").Call(Id("exclude").Op("...")),
			Return(Id("v")),
		},
		[]Code{Op("*").Id(g.names().QuerysetStruct)},
	)
}

func (g *ManagerGenerator) Generate() {
	g.AddStruct()
	g.AddConstructor()
	g.AddFilterMethod()
	g.AddExcludeMethod()
	g.AddDeleteMethod()
	g.AddInsertMethod()
	g.AddInsertInstanceMethod()
	g.AddUpdateMethod()
	g.AddAllMethod()
	g.AddGetMethod()
	g.AddNewModelMethod()
	g.AddToSQLMethod()
}
