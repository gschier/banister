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
		defineGoType = Id("v").Qual(segments[0], segments[1])
	} else {
		defineGoType = Id("v").Id(segments[0])
	}

	// Create assign statement, depending if the value is nullable
	// or not
	var assignTo *Statement
	if f.Settings().Null {
		assignTo = Op("&").Id("v")
	} else {
		assignTo = Id("v")
	}

	// Instantiate the setter
	setterDef := Id(g.names().QuerysetSetterArgStruct).Values(Dict{
		Id("field"): Lit(f.Settings().DBColumn),
		Id("value"): assignTo,
	})

	g.AddMethod(f.Settings().Name,
		[]Code{defineGoType},
		[]Code{Return(setterDef)},
		[]Code{Id(g.names().QuerysetSetterArgStruct)},
	)
}

func (g *SettersGenerator) Generate() {
	g.File.Type().Id(g.Model.Settings().Names().QuerysetSetterOptionsStruct).Struct()
	for _, f := range g.Model.Fields() {
		g.AddSetterMethod(f)
	}
}
