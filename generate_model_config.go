package banister

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
)

type ModelConfigGenerator struct {
	Model Model
	File  *File
}

func NewModelConfigGenerator(file *File, model Model) *ModelConfigGenerator {
	return &ModelConfigGenerator{Model: model, File: file}
}

func (g *ModelConfigGenerator) AddHookField(name, timing, op string) *Statement {
	comment := fmt.Sprintf("// %s sets a hook for the model that will \n"+
		"// be called %s the model is %s into the database.", name, timing, op)
	return Comment(comment).Line().Id(name).Func().Params(
		Id("m").Op("*").Id(g.Model.Settings().Names().ModelStruct),
	)
}

func (g *ModelConfigGenerator) Generate() {
	g.File.Type().Id(g.Model.Settings().Names().ConfigStruct).Struct(
		g.AddHookField("HookPreInsert", "before", "inserted").Line(),
		g.AddHookField("HookPostInsert", "after", "inserted").Line(),
		g.AddHookField("HookPreUpdate", "before", "updated").Line(),
		g.AddHookField("HookPostUpdate", "after", "updated").Line(),
		g.AddHookField("HookPreDelete", "before", "deleted").Line(),
		g.AddHookField("HookPostDelete", "after", "deleted"),
	)
}
