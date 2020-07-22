package main

import (
	"github.com/gschier/banister"
	"github.com/gschier/banister/testutil"
)

func main() {
	config := &banister.GenerateConfig{
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
