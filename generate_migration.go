package banister

import . "github.com/dave/jennifer/jen"

type MigrationGenerator struct {
	Models []Model
	File   *File
}

func NewMigrationGenerator(file *File, models []Model) *MigrationGenerator {
	return &MigrationGenerator{File: file, Models: models}
}

func (g *MigrationGenerator) AddMigrateFunc() {
	addModels := make([]Code, 0)
	for _, m := range g.Models {
		fields := []Code{Lit(m.Settings().Name)}

		for _, f := range m.Fields() {
			fieldFn := ""
			switch f.Type() {
			case Auto:
				fieldFn = "NewAutoField"
			case Char:
				// TODO
			case DateTime:
				fieldFn = "NewDateTimeField"
			case Duration:
				// TODO
			case Integer:
				fieldFn = "NewIntegerField"
			case Text:
				fieldFn = "NewTextField"
			case TextArray:
				// TODO
			case Float:
				// TODO
			case Boolean:
				fieldFn = "NewBooleanField"
			default:
				panic("Unsupported field type " + f.Type())
			}

			fieldDef := Line().Qual("github.com/gschier/banister", fieldFn).Call(Lit(f.Settings().Name))
			if f.Settings().Null {
				fieldDef.Dot("Null").Call()
			}
			if f.Settings().Default.IsValid() {
				fieldDef.Dot("Default").Call(Lit(f.Settings().Default.Value))
			}
			fields = append(fields, fieldDef)
		}

		addModels = append(addModels, Id("NewModel").Call(
			fields...,
		).Line())
	}

	g.File.Func().Id("Migrate").Params().Params().Block(
		addModels...,
	).Line()
}

func (g *MigrationGenerator) Generate() {
	g.File.ImportAlias("github.com/gschier/banister", ".")
	g.AddMigrateFunc()
}
