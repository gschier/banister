package banister_test

import (
	. "github.com/gschier/banister"
	r "github.com/stretchr/testify/require"
	"testing"
)

func TestIntegerField(t *testing.T) {
	t.Run("builds field", func(t *testing.T) {
		s := NewIntegerField("age").Build().Settings()
		r.Equal(t, "Age", s.Name)
		r.Equal(t, "age", s.DBColumn)
		r.Equal(t, false, s.Default.IsValid())
	})
}

func TestAutoField(t *testing.T) {
	t.Run("build field", func(t *testing.T) {
		s := NewAutoField("id").Build().Settings()
		r.Equal(t, true, s.PrimaryKey)
		r.Equal(t, "ID", s.Name)
		r.Equal(t, "id", s.DBColumn)
		r.Equal(t, false, s.Default.IsValid())
	})
}

func TestTextField(t *testing.T) {
	t.Run("build field", func(t *testing.T) {
		s := NewTextField("Username").Build().Settings()
		r.Equal(t, "Username", s.Name)
		r.Equal(t, "username", s.DBColumn)
		r.Equal(t, false, s.Default.IsValid())
	})

	t.Run("build field complex", func(t *testing.T) {
		s := NewTextField("Username").
			Default("foo").Null().Unique().
			Build().Settings()
		r.Equal(t, "Username", s.Name)
		r.Equal(t, "username", s.DBColumn)
		r.Equal(t, true, s.Null)
		r.Equal(t, true, s.Unique)
		r.Equal(t, true, s.Default.IsValid())
		r.Equal(t, "foo", s.Default.Value)
	})
}
