package banister

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"strings"
)

type ModelGenerator struct {
	Model Model
	File  *jen.File
}

func NewModelGenerator(file *jen.File, model Model) *ModelGenerator {
	return &ModelGenerator{Model: model, File: file}
}

func (g *ModelGenerator) names() GeneratedModelNames {
	return g.Model.Settings().Names()
}

func (g *ModelGenerator) FieldStmt(f Field) jen.Code {
	goType := fmt.Sprintf("%T", f.EmptyDefault())
	if f.Settings().Null {
		goType = "*" + goType
	}

	segments := strings.SplitN(goType, ".", 2)

	if len(segments) == 1 {
		return jen.Id(f.Settings().Name).Id(goType)
	}

	// Add import for types that require packages like time.Time
	return jen.Id(f.Settings().Name).Qual(segments[0], segments[1])
}

func (g *ModelGenerator) AddJSONMethod() {
	g.File.Comment("PrintJSON prints out a JSON string of the model for debugging")
	g.File.Func().Params(
		jen.Id("model").Op("*").Id(g.names().ModelStruct),
	).Id("PrintJSON").Params(
	// No function args
	).Params(
	// Returns nothing
	).Block(
		jen.List(jen.Id("b"), jen.Err()).Op(":=").Qual("encoding/json", "MarshalIndent").Call(
			jen.Id("model"),
			jen.Lit(""),
			jen.Lit("  "),
		),
		jen.If(jen.Parens(jen.Err().Op("!=").Nil())).Block(jen.Panic(jen.Err())),
		jen.Qual("fmt", "Printf").Call(
			jen.Lit("%T: %s"),
			jen.Id("model"),
			jen.Id("b"),
		),
	)
}

func (g *ModelGenerator) Generate() {
	// Generate the struct field definitions
	fields := make([]jen.Code, 0)
	for _, f := range g.Model.Fields() {
		fields = append(fields, g.FieldStmt(f))
	}

	// Define the struct with its fields
	g.File.Comment(
		"// " + g.names().ModelStruct + " is a database model which represents a single row from the \n" +
			"// " + g.Model.Settings().DBTable + " database table")
	g.File.Type().Id(g.names().ModelStruct).Struct(fields...)

	g.AddJSONMethod()
}
