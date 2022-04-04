package banister_test

import (
	. "github.com/gschier/banister"
	_ "github.com/gschier/banister/backends/postgres"
	"github.com/gschier/banister/testutil"
	r "github.com/stretchr/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	models := []Model{
		testutil.TestUserModel(),
		testutil.TestPostModel(),
	}

	for _, m := range models {
		m.ProvideModels(models...)
	}

	src := GenerateToString(&GenerateConfig{
		OutputDir:   "./generate/foo/bar",
		PackageName: "dummy",
		Backend:     "postgres",
		Models:      models,
	})
	r.Contains(t, src, "type User struct {")
	r.Contains(t, src, "type UserConfig struct {")
	r.Contains(t, src, "type Post struct {")
	r.Contains(t, src, "type PostConfig struct {")
}
