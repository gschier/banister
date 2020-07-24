package banister_test

import (
	. "github.com/gschier/banister"
	_ "github.com/gschier/banister/backends/postgres"
	"github.com/gschier/banister/testutil"
	r "github.com/stretchr/testify/require"
	"testing"
)

func TestBuildTableSQL(t *testing.T) {
	t.Run("builds simple table", func(t *testing.T) {
		s := BuildTableSQL(GetBackend("postgres"), testutil.TestUserModel())
		r.Equal(t, "CREATE TABLE users ( "+
			"id SERIAL PRIMARY KEY, "+
			"age INTEGER, "+
			"admin BOOLEAN NOT NULL, "+
			"name TEXT NOT NULL, "+
			"username TEXT NOT NULL UNIQUE, "+
			"created TIMESTAMP WITH TIME ZONE NOT NULL, "+
			"bio TEXT DEFAULT 'Fill me in' NOT NULL "+
			");", s)
	})
}
