package banister

import "github.com/dave/jennifer/jen"

type MigrationGenerator struct {
	Models []model
	File   *jen.File
}

func NewMigrationGenerator(file *jen.File, models []model) *MigrationGenerator {
	return &MigrationGenerator{File: file, Models: models}
}

func (g *MigrationGenerator) Generate() {
	g.File.Comment("Migration")
}
