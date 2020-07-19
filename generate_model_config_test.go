package banister_test

import (
	"github.com/dave/jennifer/jen"
	. "github.com/gschier/banister"
	"github.com/gschier/banister/testutil"
	assert "github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestModelConfigGenerator_GoSource(t *testing.T) {
	file := jen.NewFile("dummy")
	NewModelConfigGenerator(file, testutil.TestUserModel()).Generate()
	assert.Equal(t, strings.TrimSpace(`
package dummy

type UserConfig struct {
	// HookPreInsert sets a hook for the model that will
	// be called before the model is inserted into the database.
	HookPreInsert func(m *User)

	// HookPostInsert sets a hook for the model that will
	// be called after the model is inserted into the database.
	HookPostInsert func(m *User)

	// HookPreUpdate sets a hook for the model that will
	// be called before the model is updated into the database.
	HookPreUpdate func(m *User)

	// HookPostUpdate sets a hook for the model that will
	// be called after the model is updated into the database.
	HookPostUpdate func(m *User)

	// HookPreDelete sets a hook for the model that will
	// be called before the model is deleted into the database.
	HookPreDelete func(m *User)

	// HookPostDelete sets a hook for the model that will
	// be called after the model is deleted into the database.
	HookPostDelete func(m *User)
}
`), strings.TrimSpace(file.GoString()))
}
