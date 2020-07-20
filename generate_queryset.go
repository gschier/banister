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

func (g *QuerysetGenerator) AddOrderMethod() {
	g.AddChainedMethod("Order",
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
	g.AddMethod("Update",
		[]Code{Id("set").Op("...").Id(g.names().QuerysetSetterArgStruct)},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Error()},
	)
}

func (g *QuerysetGenerator) AddAllMethod() {
	g.AddMethod("All",
		[]Code{},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Index().Id(g.names().ModelStruct), Error()},
	)
}

func (g *QuerysetGenerator) AddOneMethod() {
	g.AddMethod("One",
		[]Code{},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Op("*").Id(g.names().ModelStruct), Error()},
	)
}

func (g *QuerysetGenerator) AddCountMethod() {
	g.AddMethod("Count",
		[]Code{},
		[]Code{Panic(Lit("implement me"))},
		[]Code{Int(), Error()},
	)
}

func (g *QuerysetGenerator) AddDeleteMethod() {
	defineQuery := Id("query").Op(":=").Qual("github.com/Masterminds/squirrel", "Delete").Call(
		Lit(g.Model.Settings().DBTable),
	)

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

	g.AddMethod("Delete",
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

func (g *QuerysetGenerator) AddStarSelectMethod() {
	selectArgs := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		name := f.Settings().Names(g.Model).QualifiedColumn
		selectArgs = append(selectArgs, Line().Lit(name))
	}

	// query := squirrel.Select(...).From(...)
	defineQuery := Id("query").Op(":=").
		Qual("github.com/Masterminds/squirrel", "Select").Call(selectArgs...).
		Dot("From").Call(Lit(g.Model.Settings().DBTable))

	// joinCheck := make(map[string]bool)
	defineJoinCheck := Id("joinCheck").Op(":=").
		Id("make").Call(Id("map").Index(String()).Bool())

	// Loop through things and add joins
	// for _, w := range f.filters {
	//   query = query.Where(w.filter)
	//   for _, j := range w.joins {
	//     if _, ok := joinCheck[j]; ok {
	//       continue
	//     }
	//     joinCheck[j] = true
	//     query = query.Join(j)
	//   }
	// }
	loopAndJoin := Comment("Assign filters and join if necessary").Line().For(
		Op("_").Op(",").Id("w").Op(":=").Range().Id("qs").Dot("filter"),
	).Block(
		Id("query").Op("=").Id("query").Dot("Where").Call(
			Id("w").Dot("filter"),
		),
		For(
			Op("_").Op(",").Id("j").Op(":=").Range().Id("w").Dot("joins"),
		).Block(
			If(Op("_").Op(",").Id("ok").Op(":=").
				Id("joinCheck").Index(Id("j")).Op(";").Id("ok")).
				Block(
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

	// Apply default sorts
	// if len(f.sorts) == 0 {
	//	 return query.OrderBy(
	//	`  "UsersOverride"."created" DESC`,
	//	 )
	// }
	applyDefaultOrder := Comment("Apply default order if none specified").Line().If(
		Id("len").Call(Id("qs").Dot("orderBy")).Op("==").Lit(0),
	).Block(
		Comment("TODO: Add default order-by"),
		// TODO: Add default order-by support
		//  jen.Return(...),
	)

	// // Apply user-specified sorts
	// for _, s := range f.sorts {
	//   query = query.OrderBy(s.field + " " + s.order)
	//   if s.join == "" {
	//     continue
	//   }
	//
	//   // Have we already added this join to the query?
	// 	 if _, ok := joinCheck[s.join]; ok {
	// 	   continue
	// 	 }
	//
	//   joinCheck[s.join] = true
	//   query = query.Join(s.join)
	// }
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
func (g *QuerysetGenerator) AddChainedMethod(name string, args []Code, block []Code) {
	returnType := Op("*").Id(g.names().QuerysetStruct)
	g.AddMethod(name, args, block, []Code{returnType})
}

// AddMethod is a helper to add a struct method
func (g *QuerysetGenerator) AddMethod(name string, args []Code, block []Code, returns []Code) {
	receiver := Id("qs").Op("*").Id(g.names().QuerysetStruct)
	g.File.Func().Params(receiver).Id(name).Params(args...).Params(returns...).Block(block...)
}

func (g *QuerysetGenerator) AddFilterOptionsStruct() {
	fields := make([]Code, 0)
	values := Dict{}
	for _, f := range g.Model.Fields() {
		fieldName := f.Settings().Name
		fieldType := f.Settings().Names(g.Model).FilterOptionStruct
		fields = append(fields, Id(fieldName).Id(fieldType))
		values[Id(fieldName)] = Id(fieldType).Values()
	}

	// Define options struct and instantiate an instance inline, so we
	// don't clutter the scope with unnecessary type.
	g.File.Var().Id(g.names().FilterOptionsVar).Op("=").Struct(fields...).Values(values)
}

func (g *QuerysetGenerator) Generate() {
	// Create main struct and constructor
	g.AddStruct()
	g.AddConstructor()

	// Methods
	g.AddFilterMethod()
	g.AddOrderMethod()
	g.AddLimitMethod()
	g.AddOffsetMethod()
	g.AddUpdateMethod()
	g.AddAllMethod()
	g.AddOneMethod()
	g.AddDeleteMethod()
	g.AddCountMethod()
	g.AddScanMethod()
	g.AddStarSelectMethod()
	g.AddToSQLMethod()

	// Other types
	g.AddFilterArgsStruct()
	g.AddOrderByArgsStruct()
	g.AddSetterArgsStruct()

	// Where helper
	g.AddFilterOptionsStruct()
}
