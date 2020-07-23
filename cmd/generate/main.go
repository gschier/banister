package main

import (
	"github.com/gschier/banister"
	_ "github.com/gschier/banister/backends/postgres"
	_ "github.com/gschier/banister/backends/sqlite"
	"github.com/gschier/banister/testutil"
)

func main() {
	config := &banister.GenerateConfig{
		Backend:     "sqlite3",
		OutputDir:   "./integration/generated",
		PackageName: "gen",
		MultiFile:   true,
		Models: []banister.Model{
			testutil.TestUserModel(),
			testutil.TestPostModel(),
		},
	}

	err := banister.Generate(config)
	if err != nil {
		panic(err)
	}
}
