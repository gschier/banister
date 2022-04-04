Go Banister
===========

`banister` provides access to relational databases in the safest way possible.

```go
// Example usage
package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"models"
)

func main() {
	db, _ := sql.Open("sqlite3", ":memory:?_fk=1")
	store := models.NewStore(db, models.StoreConfig{
		BlogPostConfig: models.BlogPostConfig{
			HookPreInsert: func(m *models.User) {
				m.CreatedAt = time.Now()
			},
		},
	})

	user := store.Users.MustInsert(
		models.Set.User.Username("superUser"),
		models.Set.User.Email("superuser@example.com"),
	)
	user.Email = "superuser+1@example.com"
	store.Users.MustUpdate(user)

	store.BlogPosts.MultInsert(
		models.Set.BlogPost.Title("First Post!"),
		models.Set.BlogPost.Slug("first-post"),
		models.Set.BlogPost.Published(true),
		models.Set.BlogPost.Content("Hello World!"),
		models.Set.BlogPost.UserId(user.ID),
	)

	// Fetch blog posts
	drafts := store.BlogPosts.Filter(
		models.Where.BlogPost.Published.False(),
	).Sort(
		models.OrderBy.BlogPost.CreatedAt.Desc,
	).MustAll()
	
	fmt.Printf("Found %d drafts", len(drafts))
}
```

```go
// Define models and generate a new type-safe database client
package main

import (
	"github.com/gschier/banister"
	"log"
)

var User = banister.NewModel(
	"User",
	banister.NewAutoField("id"),
	banister.NewTextField("email").Unique(),
	banister.NewTextField("password"),
)

var BlogPost = banister.NewModel(
	"BlogPost",
	banister.NewAutoField("id"),
	banister.NewDateTimeField("created_at"),
	banister.NewTextField("content"),
	banister.NewTextField("slug").Unique(),
	banister.NewTextField("title"),
	banister.NewBooleanField("published"),
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
