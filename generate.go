package banister

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
)

var __backend Backend = nil

type GenerateConfig struct {
	Backend     string
	PackageName string
	OutputDir   string
	Models      []Model
	MultiFile   bool
}

func Generate(c *GenerateConfig) error {
	files := generateJen(c)
	err := os.MkdirAll(c.OutputDir, 0755)
	if err != nil {
		return err
	}

	for name, f := range files {
		err = f.Save(filepath.Join(c.OutputDir, name))
		if err != nil {
			return err
		}
	}

	fmt.Println("Generated source to", c.OutputDir)
	return nil
}

func GenerateToString(c *GenerateConfig) string {
	for _, f := range generateJen(c) {
		// Only take the first one
		return f.GoString()
	}

	panic("No files generated")
}

func generateJen(c *GenerateConfig) map[string]*jen.File {
	__backend = GetBackend(c.Backend)

	// Initialize models
	for _, m := range c.Models {
		m.ProvideModels(c.Models...)
	}

	files := make(map[string]*jen.File)

	// Generate things per model
	NewStoreGenerator(file(c, files, globalNames.StoreStruct), c.Models).Generate()
	//NewMigrationGenerator(file(c, files, "migrations/migrations"), c.Models).Generate()
	for _, m := range c.Models {
		n := m.Settings().Names()
		NewModelGenerator(file(c, files, n.ModelStruct+"Model"), m).Generate()
		NewModelConfigGenerator(file(c, files, n.ConfigStruct), m).Generate()
		NewManagerGenerator(file(c, files, n.ManagerStruct), m).Generate()
		NewQuerysetGenerator(file(c, files, n.QuerysetStruct), m).Generate()
		NewSettersGenerator(file(c, files, n.QuerysetSetterOptionsStruct), m).Generate()
		NewOrderBysGenerator(file(c, files, n.QuerysetOrderByOptionsStruct), m).Generate()
		NewFilterGenerator(file(c, files, n.QuerysetFilterOptionsStruct), m).Generate()
	}

	return files
}

func file(c *GenerateConfig, files map[string]*jen.File, relPath string) *jen.File {
	if !c.MultiFile {
		relPath = "Gen"
	}

	relPath = strcase.ToSnake(relPath) + ".go"
	if f, ok := files[relPath]; ok {
		return f
	}

	fullPath := filepath.Join(c.OutputDir, relPath)
	dir := filepath.Dir(fullPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	f := jen.NewFilePathName(dir, c.PackageName)
	f.ImportAlias("github.com/Masterminds/squirrel", "sq")
	f.PackageComment("Code generated by banister; DO NOT EDIT.")
	files[relPath] = f
	return f
}
