package banister_test

import (
	. "github.com/gschier/banister"
	. "github.com/gschier/banister/testutil"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestMigration_RenameField(t *testing.T) {
	m := Migration{Models: []Model{TestUserModel()}}
	s := m.RenameField("user", "username", "handle")
	assert.Equal(t, "RENAME COLUMN username TO handle;", s)
}
