package banister

import . "github.com/dave/jennifer/jen"

type QuerysetGenerator struct {
	File  *File
	Model Model
}

func NewQuerysetGenerator(file *File, model Model) *QuerysetGenerator {
	return &QuerysetGenerator{File: file, Model: model}
}

func (g *QuerysetGenerator) names() GeneratedModelNames {
	return g.Model.Settings().Names()
}

func (g *QuerysetGenerator) AddFilterArgsStruct() {
	g.File.Type().Id(g.names().QuerysetFilterArgStruct).Struct(
		Id("filter").Qual("github.com/Masterminds/squirrel", "Sqlizer"),
		Id("joins").Index().String(),
	)
}

func (g *QuerysetGenerator) AddOrderByArgsStruct() {
	g.File.Type().Id(g.names().QuerysetOrderByArgStruct).Struct(
		Id("field").String(),
		Id("order").String(),
		Id("join").String(),
	)
}

func (g *QuerysetGenerator) AddSetterArgsStruct() {
	g.File.Type().Id(g.names().QuerysetSetterArgStruct).Struct(
		Id("field").String(),
		Id("value").Interface(),
	)
}

func (g *QuerysetGenerator) AddConstructor() {
	g.File.Func().Id(g.names().QuerysetConstructor).Params(
		Id("mgr").Op("*").Id(g.names().ManagerStruct),
	).Params(Op("*").Id(g.names().QuerysetStruct)).Block(
		Return(
			Op("&").Id(g.names().QuerysetStruct).Values(Dict{
				Id("mgr"):     Id("mgr"),
				Id("filter"):  Id("make").Call(Index().Id(g.names().QuerysetFilterArgStruct), Lit(0)),
				Id("orderBy"): Id("make").Call(Index().Id(g.names().QuerysetOrderByArgStruct), Lit(0)),
				Id("limit"):   Lit(0),
				Id("offset"):  Lit(0),
			}),
		),
	)
}

func (g *QuerysetGenerator) AddFilterMethod() {
	g.AddChainedMethod("Filter",
		[]Code{Id("filter").Op("...").Id(g.names().QuerysetFilterArgStruct)},
		[]Code{
			Id("qs").Dot("filter").Op("=").Id("append").Params(
				Id("qs").Dot("filter"),
				Id("filter").Op("..."),
			),
			Return(Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddSortMethod() {
	g.AddChainedMethod("Sort",
		[]Code{Id("orderBy").Op("...").Id(g.names().QuerysetOrderByArgStruct)},
		[]Code{
			Id("qs").Dot("orderBy").Op("=").Id("append").Params(
				Id("qs").Dot("orderBy"),
				Id("orderBy").Op("..."),
			),
			Return(Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddLimitMethod() {
	g.AddChainedMethod("Limit",
		[]Code{Id("limit").Uint64()},
		[]Code{
			Id("qs").Dot("limit").Op("=").Id("limit"),
			Return(Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddOffsetMethod() {
	g.AddChainedMethod("Offset",
		[]Code{Id("offset").Uint64()},
		[]Code{
			Id("qs").Dot("offset").Op("=").Id("offset"),
			Return(Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddUpdateMethod() {
	defineQuery := Id("query").Op(":=").
		Qual("github.com/Masterminds/squirrel", "Update").
		Call(Lit(g.Model.Settings().DBTable))

	mapSets := Comment("Apply setters to query").Line().For(
		Op("_").Op(",").Id("s").Op(":=").Range().Id("set"),
	).Block(
		Id("query").Op("=").Id("query").Dot("Set").Call(
			Id("s").Dot("field"),
			Id("s").Dot("value"),
		),
	)

	mapFilters := Comment("Apply filters to query").Line().For(
		Op("_").Op(",").Id("f").Op(":=").Range().Id("qs").Dot("filter"),
	).Block(
		Id("query").Op("=").Id("query").Dot("Where").Call(
			Id("f").Dot("filter"),
		),
	)

	getSQL := Id("q").Op(",").Id("args").Op(":=").
		Id("qs").Dot("toSQL").Call(Id("query"))

	exec := Id("_").Op(",").Err().Op(":=").Id("qs").Dot("mgr").Dot("db").Dot("Exec").Call(
		Id("q"),
		Id("args").Op("..."),
	)

	g.AddMethodWithPanicVariant("Update",
		[]Code{Id("set").Op("...").Id(g.names().QuerysetSetterArgStruct)},
		[]Code{Id("set").Op("...")},
		[]Code{
			defineQuery.Line(),
			mapSets.Line(),
			mapFilters.Line(),
			getSQL,
			exec,
			Return(Err()),
		},
		[]Code{Error()},
	)
}

func (g *QuerysetGenerator) AddAllMethod() {
	selectStar := Id("query").Op(":=").Id("qs").Dot("starSelect").Call()

	getSQL := Id("q").Op(",").Id("args").Op(":=").
		Id("qs").Dot("toSQL").Call(Id("query"))

	exec := Id("rows").Op(",").Err().Op(":=").Id("qs").Dot("mgr").Dot("db").Dot("Query").Call(
		Id("q"),
		Id("args").Op("..."),
	)

	checkErr := If(Err().Op("!=").Nil()).Block(
		Return(Nil(), Err()),
	)

	deferClose := Defer().Id("rows").Dot("Close").Call()

	defineResults := Id("items").Op(":=").Id("make").Call(
		Index().Id(g.names().ModelStruct),
		Lit(0),
	)

	mapResults := For(Id("rows").Dot("Next").Call()).Block(
		Var().Id("m").Id(g.names().ModelStruct),
		Err().Op("=").Id("qs").Dot("scan").Call(Id("rows"), Op("&").Id("m")),
		checkErr,
		Id("items").Op("=").Id("append").Call(Id("items"), Id("m")),
	)

	g.AddMethodWithPanicVariant("All",
		[]Code{},
		[]Code{},
		[]Code{
			selectStar,
			getSQL.Line(),
			exec,
			checkErr,
			deferClose.Line(),
			defineResults.Line(),
			mapResults.Line(),
			Return(Id("items"), Err()),
		},
		[]Code{Index().Id(g.names().ModelStruct), Error()},
	)
}

func (g *QuerysetGenerator) AddOneMethod() {
	callAll := Id("items").Op(",").Err().Op(":=").Id("qs").Dot("All").Call()

	checkLen := If(Len(Id("items")).Op("==").Lit(0)).Block(
		Return(Nil(), Qual("database/sql", "ErrNoRows")),
	)

	g.AddMethodWithPanicVariant("One",
		[]Code{},
		[]Code{},
		[]Code{
			callAll.Line(),
			Comment("Ensure we have a result"),
			checkLen.Line(),
			Return(Op("&").Id("items").Index(Lit(0)), Err()),
		},
		[]Code{Op("*").Id(g.names().ModelStruct), Error()},
	)
}

func (g *QuerysetGenerator) AddDeleteMethod() {
	defineQuery := Id("query").Op(":=").
		Qual("github.com/Masterminds/squirrel", "Delete").
		Call(Lit(g.Model.Settings().DBTable))

	applyFilters := For(
		Op("_").Op(",").Id("w").Op(":=").Range().Id("qs").Dot("filter"),
	).Block(
		Id("query").Op("=").Id("query").Dot("Where").Call(
			Id("w").Dot("filter"),
		),
	)

	toSql := Id("q").Op(",").Id("args").Op(":=").
		Id("qs").Dot("toSQL").Call(Id("query"))

	execDB := Op("_").Op(",").Err().Op(":=").
		Id("qs").Dot("mgr").Dot("db").Dot("Exec").
		Call(
			Id("q"),
			Id("args").Op("..."),
		)

	g.AddMethodWithPanicVariant("Delete",
		[]Code{},
		[]Code{},
		[]Code{
			defineQuery.Line(),
			applyFilters.Line(),
			toSql,
			execDB,
			Return(Err()),
		},
		[]Code{Error()},
	)
}

func (g *QuerysetGenerator) AddScanMethod() {
	fields := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		fields = append(fields, Op("&").Id("m").Dot(f.Settings().Name))
	}

	g.AddMethod("scan",
		[]Code{
			Id("r").Op("*").Qual("database/sql", "Rows"),
			Id("m").Op("*").Id(g.names().ModelStruct),
		},
		[]Code{Return(Id("r").Dot("Scan").Call(fields...))},
		[]Code{Error()},
	)
}

func (g *QuerysetGenerator) AddToSQLMethod() {
	g.AddMethod("toSQL",
		[]Code{Id("q").Qual("github.com/Masterminds/squirrel", "Sqlizer")},
		[]Code{
			Return(Id("qs").Dot("mgr").Dot("toSQL").Call(Id("q"))),
		},
		[]Code{
			String(),
			Index().Interface(),
		},
	)
}

func (g *QuerysetGenerator) AddStarSelectMethod() {
	selectArgs := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		name := f.Settings().Names(g.Model).QualifiedColumn
		selectArgs = append(selectArgs, Line().Lit(name))
	}

	defineQuery := Id("query").Op(":=").
		Qual("github.com/Masterminds/squirrel", "Select").Call(selectArgs...).
		Dot("From").Call(Lit(g.Model.Settings().DBTable))

	defineJoinCheck := Id("joinCheck").Op(":=").
		Id("make").Call(Id("map").Index(String()).Bool())

	loopAndJoin := Comment("Assign filters and join if necessary").Line().For(
		Op("_").Op(",").Id("w").Op(":=").Range().Id("qs").Dot("filter"),
	).Block(
		Id("query").Op("=").Id("query").Dot("Where").Call(
			Id("w").Dot("filter"),
		),
		For(
			Op("_").Op(",").Id("j").Op(":=").Range().Id("w").Dot("joins"),
		).Block(
			If(
				Op("_").Op(",").Id("ok").Op(":=").Id("joinCheck").Index(Id("j")).Op(";").Id("ok"),
			).Block(
				Continue(),
			),
			Id("joinCheck").Index(Id("j")).Op("=").True(),
			Id("query").Op("=").Id("query").Dot("Join").Call(Id("j")),
		),
	)

	// Apply limit if necessary
	applyLimit := Comment("Apply limit if set").Line().If(
		Id("qs").Dot("limit").Op(">").Lit(0),
	).Block(
		Id("query").Op("=").Id("query").Dot("Limit").Call(
			Id("qs").Dot("limit"),
		),
	)

	// Apply offset if necessary
	applyOffset := Comment("Apply offset if set").Line().If(
		Id("qs").Dot("offset").Op(">").Lit(0),
	).Block(
		Id("query").Op("=").Id("query").Dot("Offset").Call(
			Id("qs").Dot("offset"),
		),
	)

	applyDefaultOrder := Comment("Apply default order if none specified").Line().If(
		Id("len").Call(Id("qs").Dot("orderBy")).Op("==").Lit(0),
	).Block(
		Comment("TODO: Add default order-by"),
		// TODO: Add default order-by support
		//  jen.Return(...),
	)

	applyOrderBy := Comment("Apply user-specified order").Line().For(
		Op("_").Op(",").Id("s").Op(":=").Range().Id("qs").Dot("orderBy"),
	).Block(
		Id("query").Op("=").Id("query").Dot("OrderBy").Call(
			Id("s").Dot("field").Op("+").Lit(" ").Op("+").Id("s").Dot("order"),
		),
	)

	g.AddMethod("starSelect",
		[]Code{},
		[]Code{
			defineQuery.Line(),
			defineJoinCheck.Line(),
			loopAndJoin.Line(),
			applyLimit.Line(),
			applyOffset.Line(),
			applyDefaultOrder.Line(),
			applyOrderBy.Line(),
			Return(Id("query")),
		},
		[]Code{
			Qual("github.com/Masterminds/squirrel", "SelectBuilder"),
		},
	)
}

func (g *QuerysetGenerator) AddStruct() {
	g.File.Type().Id(g.names().QuerysetStruct).Struct(
		Id("mgr").Op("*").Id(g.names().ManagerStruct),
		Id("filter").Index().Id(g.names().QuerysetFilterArgStruct),
		Id("orderBy").Index().Id(g.names().QuerysetOrderByArgStruct),
		Id("limit").Uint64(),
		Id("offset").Uint64(),
	)
}

// AddChainedMethod is a helper to add a struct method that returns an instance
// of the struct for chaining.
func (g *QuerysetGenerator) AddChainedMethod(name string, args, block []Code) {
	returnType := Op("*").Id(g.names().QuerysetStruct)
	g.AddMethod(name, args, block, []Code{returnType})
}

// AddMethod is a helper to add a struct method
func (g *QuerysetGenerator) AddMethod(name string, args, block, returns []Code) {
	receiver := Id("qs").Op("*").Id(g.names().QuerysetStruct)
	g.File.Func().Params(receiver).Id(name).
		Params(args...).
		Params(returns...).
		Block(block...)
}

func (g *QuerysetGenerator) AddMethodWithPanicVariant(name string, args, panicArgs, block, returns []Code) {
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
			Err().Op(":=").Id("qs").Dot(name).Call(panicArgs...),
			If(Err().Op("!=").Nil()).Block(Panic(Err())),
		}
	} else {
		// Original method returned a value too, so handle that as well
		panicVariantReturns = []Code{returns[0]}
		callOriginalAndMaybePanic = []Code{
			Id("v").Op(",").Err().Op(":=").Id("qs").Dot(name).Call(panicArgs...),
			If(Err().Op("!=").Nil()).Block(Panic(Err())),
			Return(Id("v")),
		}
	}

	g.AddMethod(panicVariantName, args, callOriginalAndMaybePanic, panicVariantReturns)
}

func (g *QuerysetGenerator) Generate() {
	// Create main struct and constructor
	g.AddStruct()
	g.AddConstructor()

	// Methods
	g.AddFilterMethod()
	g.AddSortMethod()
	g.AddLimitMethod()
	g.AddOffsetMethod()
	g.AddUpdateMethod()
	g.AddAllMethod()
	g.AddOneMethod()
	g.AddDeleteMethod()
	g.AddScanMethod()
	g.AddStarSelectMethod()
	g.AddToSQLMethod()

	// Other types
	g.AddFilterArgsStruct()
	g.AddOrderByArgsStruct()
	g.AddSetterArgsStruct()
}
