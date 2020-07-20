package banister

import "github.com/dave/jennifer/jen"

type QuerysetGenerator struct {
	File  *jen.File
	Model Model
}

func NewQuerysetGenerator(file *jen.File, model Model) *QuerysetGenerator {
	return &QuerysetGenerator{File: file, Model: model}
}

func (g *QuerysetGenerator) names() GeneratedModelNames {
	return g.Model.Settings().Names()
}

func (g *QuerysetGenerator) AddFilterArgsStruct() {
	g.File.Type().Id(g.names().QuerysetFilterArgStruct).Struct(
		jen.Id("filter").Qual("github.com/Masterminds/squirrel", "Sqlizer"),
		jen.Id("joins").Index().String(),
	)
}

func (g *QuerysetGenerator) AddOrderByArgsStruct() {
	g.File.Type().Id(g.names().QuerysetOrderByArgStruct).Struct(
		jen.Id("field").String(),
		jen.Id("order").String(),
		jen.Id("join").String(),
	)
}

func (g *QuerysetGenerator) AddSetterArgsStruct() {
	g.File.Type().Id(g.names().QuerysetSetterArgStruct).Struct(
		jen.Id("field").String(),
		jen.Id("value").Interface(),
	)
}

func (g *QuerysetGenerator) AddConstructor() {
	g.File.Func().Id(g.names().QuerysetConstructor).Params(
		jen.Id("mgr").Op("*").Id(g.names().ManagerStruct),
	).Params(jen.Op("*").Id(g.names().QuerysetStruct)).Block(
		jen.Return(
			jen.Op("&").Id(g.names().QuerysetStruct).Values(jen.Dict{
				jen.Id("mgr"):     jen.Id("mgr"),
				jen.Id("filter"):  jen.Id("make").Call(jen.Index().Id(g.names().QuerysetFilterArgStruct), jen.Lit(0)),
				jen.Id("orderBy"): jen.Id("make").Call(jen.Index().Id(g.names().QuerysetOrderByArgStruct), jen.Lit(0)),
				jen.Id("limit"):   jen.Lit(0),
				jen.Id("offset"):  jen.Lit(0),
			}),
		),
	)
}

func (g *QuerysetGenerator) AddFilterMethod() {
	g.AddChainedMethod("Filter",
		[]jen.Code{jen.Id("filter").Op("...").Id(g.names().QuerysetFilterArgStruct)},
		[]jen.Code{
			jen.Id("qs").Dot("filter").Op("=").Id("append").Params(
				jen.Id("qs").Dot("filter"),
				jen.Id("filter").Op("..."),
			),
			jen.Return(jen.Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddOrderMethod() {
	g.AddChainedMethod("Order",
		[]jen.Code{jen.Id("orderBy").Op("...").Id(g.names().QuerysetOrderByArgStruct)},
		[]jen.Code{
			jen.Id("qs").Dot("orderBy").Op("=").Id("append").Params(
				jen.Id("qs").Dot("orderBy"),
				jen.Id("orderBy").Op("..."),
			),
			jen.Return(jen.Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddLimitMethod() {
	g.AddChainedMethod("Limit",
		[]jen.Code{jen.Id("limit").Uint64()},
		[]jen.Code{
			jen.Id("qs").Dot("limit").Op("=").Id("limit"),
			jen.Return(jen.Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddOffsetMethod() {
	g.AddChainedMethod("Offset",
		[]jen.Code{jen.Id("offset").Uint64()},
		[]jen.Code{
			jen.Id("qs").Dot("offset").Op("=").Id("offset"),
			jen.Return(jen.Id("qs")),
		},
	)
}

func (g *QuerysetGenerator) AddDeleteMethod() {
	defineQuery := jen.Id("query").Op(":=").Qual("github.com/Masterminds/squirrel", "Delete").Call(
		jen.Lit(g.Model.Settings().DBTable),
	)

	applyFilters := jen.For(
		jen.Op("_").Op(",").Id("w").Op(":=").Range().Id("qs").Dot("filter"),
	).Block(
		jen.Id("query").Op("=").Id("query").Dot("Where").Call(jen.Id("w").Dot("filter")),
	)

	toSql := jen.Id("q").Op(",").Id("args").Op(":=").Id("qs").Dot("toSQL").Call(jen.Id("query"))
	execDB := jen.Op("_").Op(",").Err().Op(":=").Id("qs").Dot("mgr").Dot("db").Dot("Exec").Call(
		jen.Id("q"),
		jen.Id("args").Op("..."),
	)

	g.AddMethod("Delete",
		[]jen.Code{jen.Id("m").Op("*").Id(g.names().ModelStruct)},
		[]jen.Code{
			defineQuery.Line(),
			applyFilters.Line(),
			toSql,
			execDB,
			jen.Return(jen.Err()),
		},
		[]jen.Code{jen.Error()},
	)
}

func (g *QuerysetGenerator) AddScanMethod() {
	fields := make([]jen.Code, 0)
	for _, f := range g.Model.Fields() {
		fields = append(fields, jen.Op("&").Id("m").Dot(f.Settings().Name))
	}

	g.AddMethod("scan",
		[]jen.Code{
			jen.Id("r").Op("*").Qual("database/sql", "Rows"),
			jen.Id("m").Op("*").Id(g.names().ModelStruct),
		},
		[]jen.Code{jen.Return(jen.Id("r").Dot("Scan").Call(fields...))},
		[]jen.Code{jen.Error()},
	)
}

func (g *QuerysetGenerator) AddToSQLMethod() {
	genSQL := jen.Id("query").Op(",").Id("args").Op(",").Err().Op(":=").
		Id("q").Dot("ToSql").Call()

	maybePanic := jen.If(
		jen.Err().Op("!=").Nil(),
	).Block(
		jen.Panic(jen.Err()),
	)

	g.AddMethod("toSQL",
		[]jen.Code{jen.Id("q").Qual("github.com/Masterminds/squirrel", "Sqlizer")},
		[]jen.Code{
			genSQL,
			maybePanic,
			jen.Return(jen.Id("query"), jen.Id("args")),
		},
		[]jen.Code{
			jen.String(),
			jen.Index().Interface(),
		},
	)
}

func (g *QuerysetGenerator) AddStarSelectMethod() {
	selectArgs := make([]jen.Code, 0)
	for _, f := range g.Model.Fields() {
		name := f.Settings().Names(g.Model).QualifiedColumn
		selectArgs = append(selectArgs, jen.Line().Lit(name))
	}

	// query := squirrel.Select(...).From(...)
	defineQuery := jen.Id("query").Op(":=").Qual("github.com/Masterminds/squirrel", "Select").Call(
		selectArgs...,
	).Dot("From").Call(
		jen.Lit(g.Model.Settings().DBTable),
	)

	// joinCheck := make(map[string]bool)
	defineJoinCheck := jen.Id("joinCheck").Op(":=").Id("make").Call(
		jen.Id("map").Index(jen.String()).Bool(),
	)

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
	loopAndJoin := jen.Comment("Assign filters and join if necessary").Line().For(
		jen.Op("_").Op(",").Id("w").Op(":=").Range().Id("qs").Dot("filter"),
	).Block(
		jen.Id("query").Op("=").Id("query").Dot("Where").Call(jen.Id("w").Dot("filter")),
		jen.For(
			jen.Op("_").Op(",").Id("j").Op(":=").Range().Id("w").Dot("joins"),
		).Block(
			jen.If(jen.Op("_").Op(",").Id("ok").Op(":=").Id("joinCheck").Index(jen.Id("j")).Op(";").Id("ok")).Block(
				jen.Continue(),
			),
			jen.Id("joinCheck").Index(jen.Id("j")).Op("=").True(),
			jen.Id("query").Op("=").Id("query").Dot("Join").Call(jen.Id("j")),
		),
	)

	// Apply limit if necessary
	applyLimit := jen.Comment("Apply limit if set").Line().If(
		jen.Id("qs").Dot("limit").Op(">").Lit(0),
	).Block(
		jen.Id("query").Op("=").Id("query").Dot("Limit").Call(
			jen.Id("qs").Dot("limit"),
		),
	)

	// Apply offset if necessary
	applyOffset := jen.Comment("Apply offset if set").Line().If(
		jen.Id("qs").Dot("offset").Op(">").Lit(0),
	).Block(
		jen.Id("query").Op("=").Id("query").Dot("Offset").Call(
			jen.Id("qs").Dot("offset"),
		),
	)

	// Apply default sorts
	// if len(f.sorts) == 0 {
	//	 return query.OrderBy(
	//	`  "UsersOverride"."created" DESC`,
	//	 )
	// }
	applyDefaultOrder := jen.Comment("Apply default order if none specified").Line().If(
		jen.Id("len").Call(
			jen.Id("qs").Dot("orderBy"),
		).Op("==").Lit(0),
	).Block(
		jen.Comment("TODO: Add default order-by"),
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
	applyOrderBy := jen.Comment("Apply user-specified order").Line().For(
		jen.Op("_").Op(",").Id("s").Op(":=").Range().Id("qs").Dot("orderBy"),
	).Block(
		jen.Id("query").Op("=").Id("query").Dot("OrderBy").Call(
			jen.Id("s").Dot("field").Op("+").Lit(" ").Op("+").Id("s").Dot("order"),
		),
	)

	g.AddMethod("starSelect",
		[]jen.Code{},
		[]jen.Code{
			defineQuery.Line(),
			defineJoinCheck.Line(),
			loopAndJoin.Line(),
			applyLimit.Line(),
			applyOffset.Line(),
			applyDefaultOrder.Line(),
			applyOrderBy.Line(),
			jen.Return(jen.Id("query")),
		},
		[]jen.Code{jen.Qual("github.com/Masterminds/squirrel", "SelectBuilder")},
	)
}

func (g *QuerysetGenerator) AddStruct() {
	g.File.Type().Id(g.names().QuerysetStruct).Struct(
		jen.Id("mgr").Op("*").Id(g.names().ManagerStruct),
		jen.Id("filter").Index().Id(g.names().QuerysetFilterArgStruct),
		jen.Id("orderBy").Index().Id(g.names().QuerysetOrderByArgStruct),
		jen.Id("limit").Uint64(),
		jen.Id("offset").Uint64(),
	)
}

// AddChainedMethod is a helper to add a struct method that returns an instance
// of the struct for chaining.
func (g *QuerysetGenerator) AddChainedMethod(name string, args []jen.Code, block []jen.Code) {
	returnType := jen.Op("*").Id(g.names().QuerysetStruct)
	g.AddMethod(name, args, block, []jen.Code{returnType})
}

// AddMethod is a helper to add a struct method
func (g *QuerysetGenerator) AddMethod(name string, args []jen.Code, block []jen.Code, returns []jen.Code) {
	receiver := jen.Id("qs").Op("*").Id(g.names().QuerysetStruct)
	g.File.Func().Params(receiver).Id(name).Params(args...).Params(returns...).Block(block...)
}

func (g *QuerysetGenerator) AddFilterOptionsStruct() {
	fields := make([]jen.Code, 0)
	values := jen.Dict{}
	for _, f := range g.Model.Fields() {
		fieldName := f.Settings().Name
		fieldType := f.Settings().Names(g.Model).FilterOptionStruct
		fields = append(fields, jen.Id(fieldName).Id(fieldType))
		values[jen.Id(fieldName)] = jen.Id(fieldType).Values()
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
	// TODO: g.AddUpdateMethod()
	// TODO: g.AddAllMethod()
	// TODO: g.AddOneMethod()
	g.AddDeleteMethod()
	// TODO: g.AddCountMethod()
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
