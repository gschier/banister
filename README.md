Go Banister
===========

`banister` provides access to relational databases in the safest way possible.

```go
package main

import (
	"github.com/gschier/banister"
	"log"
)

var User = banister.NewModel(
	"User",
	banister.NewCharField("id", 25).PrimaryKey(),
	banister.NewTextField("email").Unique(),
	banister.NewTextField("name"),
	banister.NewTextField("password"),
)

var BlogPost = banister.NewModel(
	"BlogPost",
	banister.NewCharField("id", 25).PrimaryKey(),
	banister.NewTextField("content"),
	banister.NewTextField("slug").Unique(),
	banister.NewTextField("title"),
	banister.NewForeignKeyField(User.Settings().Name).OnDelete(banister.OnDeleteSetNull),
)

func main() {
	config := &banister.GenerateConfig{
		Backend:     "postgres",
		OutputDir:   "./models",
		PackageName: "models",
		Models: []banister.Model{
			models.User,
			models.BlogPost,
		},
	}

	err := banister.Generate(config)
	if err != nil {
		log.Panicf("Failed to generate database client: %s", err)
	}
}
```
