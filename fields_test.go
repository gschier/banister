package banister_test

import (
	. "github.com/gschier/banister"
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestIntegerField(t *testing.T) {
	t.Run("builds field", func(t *testing.T) {
		s := NewIntegerField("age").Build(nil).Settings()
		assert.Equal(t, "Age", s.Name)
		assert.Equal(t, "age", s.DBColumn)
		assert.Equal(t, false, s.Default.IsValid())
	})
}

func TestAutoField(t *testing.T) {
	t.Run("build field", func(t *testing.T) {
		s := NewAutoField("id").Build(nil).Settings()
		assert.Equal(t, true, s.PrimaryKey)
		assert.Equal(t, "ID", s.Name)
		assert.Equal(t, "id", s.DBColumn)
		assert.Equal(t, false, s.Default.IsValid())
	})
}

func TestTextField(t *testing.T) {
	t.Run("build field", func(t *testing.T) {
		s := NewTextField("Username").Build(nil).Settings()
		assert.Equal(t, "Username", s.Name)
		assert.Equal(t, "username", s.DBColumn)
		assert.Equal(t, false, s.Default.IsValid())
	})

	t.Run("build field complex", func(t *testing.T) {
		s := NewTextField("Username").
			Default("foo").Null().Unique().Hidden().
			Build().Settings()
		assert.Equal(t, "Username", s.Name)
		assert.Equal(t, "username", s.DBColumn)
		assert.Equal(t, true, s.Null)
		assert.Equal(t, true, s.Unique)
		assert.Equal(t, true, s.Default.IsValid())
		assert.Equal(t, "foo", s.Default.Value)
	})
}
