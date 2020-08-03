package banister

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"strings"
	"time"
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

func (g *ModelGenerator) GenField(f Field) Code {
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

	field.Tag(map[string]string{"json": f.Settings().JSONName})
	field.Comment(BuildColumnSQL(__backend, f, true))

	return field
}

func (g *ModelGenerator) AddJSONMethod() {
	marshall := List(Id("b"), Err()).Op(":=").Qual("encoding/json", "MarshalIndent").Call(
		Id("model"),
		Lit("// "),
		Lit("  "),
	)

	checkErr := If(Parens(Err().Op("!=").Nil())).Block(Panic(Err()))

	printIt := Qual("fmt", "Printf").Call(
		Lit("\n// var %s %T = %s\n\n"),
		Lit(PrivateGoName(g.Model.Settings().Name)),
		Id("model"),
		Id("b"),
	)

	g.File.Comment("PrintJSON prints out a JSON string of the model for debugging")
	g.File.Func().Params(
		Id("model").Op("*").Id(g.names().ModelStruct),
	).Id("PrintJSON").Params( /* Args */ ).Params( /* Returns */ ).Block(
		marshall,
		checkErr.Line(),
		printIt,
	)
}

func (g *ModelGenerator) AddConstructor() {
	defaultValues := Dict{}

	for _, f := range g.Model.Fields() {
		fieldDefault := f.Settings().Default
		goDefaultVal := fieldDefault.Value

		// If no default is provided and nil is not allowed, use the
		// fallback default
		if f.Settings().Null == false && !fieldDefault.IsValid() {
			goDefaultVal = f.EmptyDefault()
		}

		var defaultValue *Statement
		switch v := goDefaultVal.(type) {
		case nil:
			// NOTE: Special case for nil defaults
			defaultValues[Id(f.Settings().Name)] = Nil()
			continue
		case time.Time:
			defaultValue = Qual("time", "Time").Values()
		case time.Duration:
			defaultValue = Qual("time", "Duration").Call()
		case []string:
			defaultValue = Make(Index().String(), Lit(0))
		default:
			defaultValue = Lit(v)
		}

		if f.Settings().Null {
			goType := Id(fmt.Sprintf("%T", f.EmptyDefault()))
			defaultValues[Id(f.Settings().Name)] = Func().Params(
				Id("v").Add(goType),
			).Params(
				Op("*").Add(goType), // Returns pointer to type
			).Block(
				Return(Op("&").Id("v")),
			).Call(defaultValue)
		} else {
			defaultValues[Id(f.Settings().Name)] = defaultValue
		}
	}

	instantiateWithDefaults := Op("&").Id(g.names().ModelStruct).Values(defaultValues)

	g.File.Comment(g.names().ModelConstructor + " returns a new instance of " +
		g.names().ModelStruct + " with default values")
	g.File.Func().Id(g.names().ModelConstructor).Params( /* Args */ ).Params(
		Op("*").Id(g.Model.Settings().Name),
	).Block(
		Return(instantiateWithDefaults),
	)
}

func (g *ModelGenerator) AddStruct() {
	// Generate the struct field definitions
	fields := make([]Code, 0)
	for _, f := range g.Model.Fields() {
		fields = append(fields, g.GenField(f))
	}

	name := g.names().ModelStruct
	g.File.Comment(name + " represents a row in the \"" + g.Model.Settings().DBTable + "\" table")
	g.File.Type().Id(g.names().ModelStruct).Struct(fields...)
}

func (g *ModelGenerator) Generate() {
	g.AddStruct()
	g.AddConstructor()
	g.AddJSONMethod()
}
