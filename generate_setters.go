package banister

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"strings"
)

type SettersGenerator struct {
	Model Model
	File  *File
}

func NewSettersGenerator(file *File, model Model) *SettersGenerator {
	return &SettersGenerator{Model: model, File: file}
}

func (g *SettersGenerator) names() GeneratedModelNames {
	return g.Model.Settings().Names()
}

// AddMethod is a helper to add a struct method
func (g *SettersGenerator) AddMethod(name string, args, block, returns []Code) {
	receiver := Id(g.names().QuerysetSetterOptionsStruct)
	g.File.Func().Params(receiver).Id(name).
		Params(args...).
		Params(returns...).
		Block(block...)
}

func (g *SettersGenerator) AddSetterMethod(f Field) {
	// If type comes from package, we need to qualify it
	goType := fmt.Sprintf("%T", f.EmptyDefault())
	segments := strings.SplitN(goType, ".", 2)
	var defineGoType *Statement
	if len(segments) == 2 {
		defineGoType = Qual(segments[0], segments[1])
	} else {
		defineGoType = Id(segments[0])
	}

	// Create assign statement, depending if the value is nullable or not
	var valueDef *Statement

	if f.Settings().Null {
		valueDef = Op("&").Id("v")
	} else {
		valueDef = Id("v")
	}

	col := fmt.Sprintf(`"%s"`, f.Settings().DBColumn)
	g.File.Comment(f.Settings().Name + " sets the " + col + " field to the provided value.")
	g.AddMethod(f.Settings().Name,
		[]Code{Id("v").Add(defineGoType)},
		[]Code{
			Return(Id(g.names().QuerysetSetterArgStruct).Values(Dict{
				Id("field"): Lit(f.Settings().DBColumn),
				Id("value"): valueDef,
			})),
		},
		[]Code{Id(g.names().QuerysetSetterArgStruct)},
	)

	// Add option to set ptr if nullable field
	if f.Settings().Null {
		g.File.Comment("// " + f.Settings().Name + "Ptr sets the " + col +
			" field but uses a pointer because\n // this field allows NULL values.")
		g.AddMethod(f.Settings().Name+"Ptr",
			[]Code{Id("v").Op("*").Add(defineGoType)},
			[]Code{
				Return(Id(g.names().QuerysetSetterArgStruct).Values(Dict{
					Id("field"): Lit(f.Settings().DBColumn),
					Id("value"): Id("v"),
				})),
			},
			[]Code{Id(g.names().QuerysetSetterArgStruct)},
		)

		// Add no-arg method to set to Null
		g.File.Comment(f.Settings().Name + "Null is helper to set the " +
			col + " field to NULL.")
		g.AddMethod(f.Settings().Name+"Null",
			[]Code{},
			[]Code{
				Return(Id(g.names().QuerysetSetterArgStruct).Values(Dict{
					Id("field"): Lit(f.Settings().DBColumn),
					Id("value"): Nil(),
				})),
			},
			[]Code{Id(g.names().QuerysetSetterArgStruct)},
		)
	}
}

func (g *SettersGenerator) Generate() {
	g.File.Type().Id(g.Model.Settings().Names().QuerysetSetterOptionsStruct).Struct()
	for _, f := range g.Model.Fields() {
		g.AddSetterMethod(f)
	}
}
