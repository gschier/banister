package integration_test

import (
	"database/sql"
	"github.com/gschier/banister"
	_ "github.com/gschier/banister/backends/sqlite"
	. "github.com/gschier/banister/integration/generated"
	"github.com/gschier/banister/testutil"
	_ "github.com/mattn/go-sqlite3"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrate(t *testing.T) {
	t.Run("Inserts a record", func(t *testing.T) {
		store, _ := createStore(t)
		user := store.Users.InsertP(Set.User.Username("another"))
		assert.Equal(t, "another", user.Username)
		assert.Equal(t, "another", user.Name, "hook should have worked")
	})

	t.Run("Filters and returns results", func(t *testing.T) {
		store, _ := createStore(t)
		_ = store.Users.InsertP(Set.User.Username("kid"), Set.User.Age(11))
		_ = store.Users.InsertP(Set.User.Username("adult"), Set.User.Age(28))

		users := store.Users.
			Filter(Where.User.Age.Gt(100)).
			Sort(OrderBy.User.Created.Asc).
			AllP()
		assert.Equal(t, 0, len(users))

		users = store.Users.
			Filter(Where.User.Age.Gt(20)).
			Sort(OrderBy.User.Created.Desc).
			AllP()
		assert.Equal(t, 2, len(users))

		users = store.Users.AllP()
		assert.Equal(t, 4, len(users))
	})

	t.Run("Deletes results", func(t *testing.T) {
		store, _ := createStore(t)
		user := store.Users.InsertP(Set.User.Username("foo"), Set.User.Age(11))

		store.Users.DeleteP(user)

		users := store.Users.Filter(Where.User.Username.Eq("foo")).AllP()
		assert.Equal(t, 0, len(users))
	})

	t.Run("Updates results", func(t *testing.T) {
		store, user := createStore(t)

		user.Name = "Baby"
		store.Users.UpdateP(user)

		users := store.Users.AllP()
		assert.Equal(t, 2, len(users))
		assert.Equal(t, "Baby", user.Name)
		assert.Equal(t, "Baby", users[0].Name, "should have updated in DB")
	})
}

func createStore(t *testing.T) (*Store, *User) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.Nil(t, err, "sqlite should open connection")

	_, err = db.Exec(banister.BuildTableSQL(
		banister.GetBackend("sqlite3"),
		testutil.TestUserModel(),
	))
	assert.Nil(t, err, "tables should be created")
	store := NewStore(db, StoreConfig{
		UserConfig: UserConfig{
			HookPreInsert: func(m *User) {
				if m.Name == "" {
					m.Name = m.Username
				}
			},
		},
	})

	// Insert some dummy data
	user := store.Users.InsertP(Set.User.Username("gschier"), Set.User.Age(11))
	store.Users.InsertP(Set.User.Username("pupper"), Set.User.Age(21))

	return store, user
}
