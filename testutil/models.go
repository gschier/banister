package testutil

import (
	"github.com/gschier/banister"
)

type TestUser struct{}

func TestUserModel() banister.Model {
	return banister.NewModel("User",
		banister.NewAutoField("id"),
		banister.NewIntegerField("age").Null(),
		banister.NewBooleanField("admin"),
		banister.NewTextField("name"),
		banister.NewTextField("username").Unique(),
		banister.NewDateTimeField("created"),
		banister.NewTextField("bio").Default("Fill me in"),
	)
}

func TestPostModel() banister.Model {
	return banister.NewModel("Post",
		banister.NewAutoField("id"),
		banister.NewForeignKeyField("User"),
		banister.NewIntegerField("score").Null().Default(50),
		banister.NewTextField("title").Default("Change Me"),
		banister.NewTextField("slug").Unique(),
		banister.NewTextField("subtitle").Null(),
		banister.NewTextField("content").Default("Fill 'this' \"in\""),
		banister.NewDateTimeField("created"),
		banister.NewBooleanField("deleted"),
		banister.NewBooleanField("approved").Null(),
		banister.NewBooleanField("private").Default(true),
	)
}
