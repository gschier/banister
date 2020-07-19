package banister_test

import (
	. "github.com/gschier/banister"
	_ "github.com/gschier/banister/backends/postgres"
	"github.com/gschier/banister/testutil"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestBuildTableSQL(t *testing.T) {
	t.Run("builds simple table", func(t *testing.T) {
		s := BuildTableSQL(GetBackend("postgres"), testutil.TestUserModel())
		assert.Equal(t, "CREATE TABLE users ( "+
			"id SERIAL NOT NULL PRIMARY KEY, "+
			"age INTEGER NULL, "+
			"name TEXT NOT NULL, " +
			"username TEXT NOT NULL UNIQUE, " +
			"created TIMESTAMP WITH TIME ZONE NOT NULL );", s)
	})
}
