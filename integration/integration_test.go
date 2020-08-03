package integration_test

import (
	"database/sql"
	"github.com/gschier/banister"
	_ "github.com/gschier/banister/backends/sqlite"
	. "github.com/gschier/banister/integration/generated"
	"github.com/gschier/banister/testutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	r "github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrate(t *testing.T) {
	t.Run("Inserts a record", func(t *testing.T) {
		store, _ := createStore(t)
		user := store.Users.InsertP(Set.User.Username("another"))
		r.Equal(t, "another", user.Username)
		r.Equal(t, "another", user.Name, "hook should have worked")
	})

	t.Run("Filters and returns results", func(t *testing.T) {
		store, _ := createStore(t)
		store.Users.InsertP(Set.User.Username("kid"), Set.User.Age(11))
		store.Users.InsertP(Set.User.Username("adult"), Set.User.Age(28))

		users := store.Users.
			Filter(Where.User.Age.Gt(100)).
			Sort(OrderBy.User.Created.Asc).
			AllP()
		r.Equal(t, 0, len(users))

		users = store.Users.
			Filter(Where.User.Age.Gt(20)).
			Sort(OrderBy.User.Created.Desc).
			AllP()
		r.Equal(t, 2, len(users))

		users = store.Users.AllP()
		r.Equal(t, 4, len(users))
	})

	t.Run("Get and One", func(t *testing.T) {
		store, _ := createStore(t)
		kid := store.Users.InsertP(Set.User.Username("kid"), Set.User.Age(11))
		store.Users.InsertP(Set.User.Username("adult"), Set.User.Age(28))

		user := store.Users.
			Filter(Where.User.Username.Eq("kid")).
			OneP()
		r.Equal(t, kid, user)

		user = store.Users.GetP(kid.ID)
		r.Equal(t, kid, user)
	})

	t.Run("Or's and And's", func(t *testing.T) {
		store, _ := createStore(t)
		store.Users.InsertP(Set.User.Username("kid"), Set.User.Age(11))
		store.Users.InsertP(Set.User.Username("admin"), Set.User.Admin(true))
		store.Users.InsertP(Set.User.Username("senior"), Set.User.Age(95))

		users := store.Users.Filter(
			Where.User.Or(
				Where.User.Admin.Eq(true),
				Where.User.And(
					Where.User.Age.Gt(90),
					Where.User.Age.Lt(100),
				),
			),
		).AllP()

		r.Equal(t, 2, len(users))
	})

	t.Run("Deletes results", func(t *testing.T) {
		store, _ := createStore(t)
		user := store.Users.InsertP(Set.User.Username("foo"), Set.User.Age(11))

		store.Users.DeleteP(user)

		users := store.Users.Filter(Where.User.Username.Eq("foo")).AllP()
		r.Equal(t, 0, len(users))
	})

	t.Run("Updates single", func(t *testing.T) {
		store, user := createStore(t)

		user.Name = "Baby"
		store.Users.UpdateP(user)

		users := store.Users.AllP()
		r.Equal(t, 2, len(users))
		r.Equal(t, "Baby", user.Name)
		r.Equal(t, "Baby", users[0].Name, "should have updated in DB")
	})

	t.Run("Updates multiple", func(t *testing.T) {
		store, _ := createStore(t)

		store.Users.Filter(
			Where.User.Username.Eq("bobby"),
		).UpdateP(
			Set.User.AgeNull(),
			Set.User.Admin(true),
		)

		users := store.Users.AllP()
		r.Equal(t, "bobby", users[0].Username)
		r.Equal(t, true, users[0].Admin)
		r.Nil(t, users[0].Age)

		r.Equal(t, false, users[1].Admin)
		r.EqualValues(t, 21, *users[1].Age)
	})

	t.Run("Bulk delete", func(t *testing.T) {
		store, _ := createStore(t)
		store.Users.Filter(Where.User.Username.Eq("bobby")).DeleteP()
		users := store.Users.AllP()
		r.Equal(t, 1, len(users))
		r.Equal(t, "tammy", users[0].Username)
	})

	t.Run("Nullable fields with defaults work", func(t *testing.T) {
		store, user := createStore(t)
		p := store.Posts.InsertP(Set.Post.UserID(user.ID))
		r.Equal(t, int64(50), *p.Score)
	})

	t.Run("Exclude query", func(t *testing.T) {
		store, user := createStore(t)
		users := store.Users.Exclude(Where.User.ID.Eq(user.ID)).AllP()
		r.Equal(t, 1, len(users))
		r.NotEqual(t, user.ID, users[0].ID)
	})

	t.Run("Text filters", func(t *testing.T) {
		store, _ := createStore(t)
		users := store.Users.Filter(Where.User.Username.Contains("obb")).AllP()
		r.Equal(t, 1, len(users))
	})

	t.Run("Complex exclude query", func(t *testing.T) {
		store, user := createStore(t)
		users := store.Users.
			Filter(
				Where.User.Age.Gt(0),
				Where.User.Age.NotEq(1123),
			).
			Exclude(Where.User.Age.IsNull()).
			Exclude(
				Where.User.Age.Eq(*user.Age),
				Where.User.Or(
					Where.User.ID.Eq(user.ID),
					Where.User.Username.Eq(user.Username),
				),
			).AllP()
		r.Equal(t, 1, len(users))
		r.NotEqual(t, user.ID, users[0].ID)
	})
}

func createStore(t *testing.T) (*Store, *User) {
	db, err := sql.Open("sqlite3", ":memory:?_fk=1")
	r.Nil(t, err, "sqlite should open connection")

	models := []banister.Model{
		testutil.TestUserModel(),
		testutil.TestPostModel(),
	}

	// NOTE: Hack to initialize models, which will not be necessary once we
	//   have generated migrations
	for _, m := range models {
		m.ProvideModels(models)
	}

	backend := banister.GetBackend("sqlite3")
	sqlStr := "" +
		banister.BuildTableSQL(backend, models[0]) + "\n" +
		banister.BuildTableSQL(backend, models[1])
	//println(sqlStr)

	_, err = db.Exec(sqlStr)
	r.Nil(t, err, "tables should be created")

	store := NewStore(db, StoreConfig{
		UserConfig: UserConfig{
			HookPreInsert: func(m *User) {
				if m.Name == "" {
					m.Name = m.Username
				}
			},
		},
	})

	user1 := store.Users.InsertP(Set.User.Username("bobby"), Set.User.Age(11))
	assert.EqualValues(t, user1.ID, 1)
	user2 := store.Users.InsertP(Set.User.Username("tammy"), Set.User.Age(21))
	assert.EqualValues(t, user2.ID, 2)

	return store, user1
}
