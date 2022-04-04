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
		model := testutil.TestPostModel()
		model.ProvideModels(testutil.TestUserModel())
		s := BuildTableSQL(GetBackend("postgres"), model)
		r.Equal(t, "CREATE TABLE posts ( "+
			"id SERIAL PRIMARY KEY, "+
			"user_id INTEGER NOT NULL, "+
			"score INTEGER DEFAULT 50, "+
			"title TEXT DEFAULT 'Change Me' NOT NULL, "+
			"slug TEXT NOT NULL UNIQUE, "+
			"subtitle TEXT, "+
			"content TEXT DEFAULT 'Fill ''this'' \"in\"' NOT NULL, "+
			"created TIMESTAMP WITH TIME ZONE NOT NULL, "+
			"deleted BOOLEAN NOT NULL, "+
			"approved BOOLEAN, "+
			"private BOOLEAN DEFAULT TRUE NOT NULL, "+
			"FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE "+
			");", s)
	})
}
