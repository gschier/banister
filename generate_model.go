package banister

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"strings"
)

type ModelGenerator struct {
	Model Model
	File  *File
}

func NewModelGenerator(file *File, model Model) *ModelGenerator {
	return &ModelGenerator{Model: model, File: file}
}

func (g *ModelGenerator) names() GeneratedModelNames {
	return g.Model.Settings().Names()
}

func (g *ModelGenerator) FieldStmt(f Field) Code {
	goType := fmt.Sprintf("%T", f.EmptyDefault())
	if f.Settings().Null {
		goType = "*" + goType
	}

	segments := strings.SplitN(goType, ".", 2)

	if len(segments) == 1 {
		return Id(f.Settings().Name).Id(goType)
	}

	// Add import for types that require packages like time.Time
	return Id(f.Settings().Name).Qual(segments[0], segments[1])
}

func (g *ModelGenerator) AddJSONMethod() {
	g.File.Comment("PrintJSON prints out a JSON string of the model for debugging")
	g.File.Func().Params(
		Id("model").Op("*").Id(g.names().ModelStruct),
	).Id("PrintJSON").Params(
	// No function args
	).Params(
	// Returns nothing
	).Block(
		List(Id("b"), Err()).Op(":=").Qual("encoding/json", "MarshalIndent").Call(
			Id("model"),
			Lit(""),
			Lit("  "),
		),
		If(Parens(Err().Op("!=").Nil())).Block(Panic(Err())),
		Qual("fmt", "Printf").Call(
			Lit("%T: %s"),
			Id("model"),
			Id("b"),
		),
	)
}

func (g *ModelGenerator) Generate() {
	// Generate the struct field definitions
	fields := make([]Code, 0)
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
