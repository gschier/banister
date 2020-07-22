package testutil

import (
	"github.com/gschier/banister"
)

type TestUser struct{}

func TestUserModel() TestUser {
	return TestUser{}
}

func (u TestUser) Settings() banister.ModelSettings {
	s := &banister.ModelSettings{Name: "User"}
	s.FillDefaults()
	return *s
}

func (u TestUser) Fields() []banister.Field {
	return []banister.Field{
		banister.NewAutoField("id").Build(),
		banister.NewIntegerField("age").Null().Build(),
		banister.NewTextField("name").Build(),
		banister.NewTextField("username").Unique().Build(),
		banister.NewDateTimeField("created").Build(),
	}
}

type TestPost struct{}

func TestPostModel() TestPost {
	return TestPost{}
}

func (u TestPost) Settings() banister.ModelSettings {
	s := &banister.ModelSettings{Name: "Post"}
	s.FillDefaults()
	return *s
}

func (u TestPost) Fields() []banister.Field {
	return []banister.Field{
		banister.NewAutoField("id").Build(),
		banister.NewIntegerField("words").Build(),
		banister.NewTextField("title").Default("Change Me").Build(),
		banister.NewTextField("slug").Unique().Build(),
		banister.NewTextField("subtitle").Null().Build(),
		banister.NewTextField("content").Default("Fill me in").Build(),
		banister.NewDateTimeField("created").Build(),
		banister.NewBooleanField("deleted").Build(),
		banister.NewBooleanField("approved").Null().Build(),
		banister.NewBooleanField("private").Default(true).Build(),
	}
}
