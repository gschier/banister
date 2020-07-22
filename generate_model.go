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

	// Add import for types that require packages like time.Time
	var field *Statement
	if len(segments) == 2 {
		field = Id(f.Settings().Name).Qual(segments[0], segments[1])
	} else {
		field = Id(f.Settings().Name).Id(goType)
	}

	return field.Tag(map[string]string{"json": f.Settings().JSONName})
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
			Lit("// "),
			Lit("  "),
		),
		If(Parens(Err().Op("!=").Nil())).Block(Panic(Err())),
		Qual("fmt", "Printf").Call(
			Lit("\n// var %s %T = %s\n\n"),
			Lit(PrivateGoName(g.Model.Settings().Name)),
			Id("model"),
			Id("b"),
		),
	)
}

func (g *ModelGenerator) AddStruct() {
	// Generate the struct field definitions
	fields := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		fields = append(fields, g.FieldStmt(f))
	}

	name := g.names().ModelStruct
	comment := "// " + name + " is a database model which represents a single row from the \n" +
		"// " + g.Model.Settings().DBTable + " database table"

	g.File.Comment(comment).Line().Type().Id(g.names().ModelStruct).Struct(fields...)
}

func (g *ModelGenerator) Generate() {
	g.AddStruct()
	g.AddJSONMethod()
}
