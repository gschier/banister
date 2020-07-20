package banister_test

import (
	. "github.com/gschier/banister"
	"github.com/gschier/banister/testutil"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	src := GenerateToString(&GenerateConfig{
		OutputDir:   "./generate/foo/bar",
		PackageName: "dummy",
		Models: []Model{
			testutil.TestUserModel(),
			testutil.TestPostModel(),
		},
	})
	assert.Contains(t, src, "type User struct {")
	assert.Contains(t, src, "type UserConfig struct {")
	assert.Contains(t, src, "type Post struct {")
	assert.Contains(t, src, "type PostConfig struct {")
}
