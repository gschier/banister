package banister

import (
	r "github.com/stretchr/testify/require"
	"testing"
)

func TestPublicGoName(t *testing.T) {
	names := map[string]string{
		"Same":          "Same",
		"small":         "Small",
		"UPPER":         "UPPER",
		"UsersALL_CAPS": "UsersALLCAPS",
		"snake_case":    "SnakeCase",
		"camelCase":     "CamelCase",
	}

	for in, out := range names {
		t.Run("should convert name "+in, func(t *testing.T) {
			r.Equal(t, out, PublicGoName(in))
		})
	}
}

func TestPrivateGoName(t *testing.T) {
	names := map[string]string{
		"Same":          "same",
		"small":         "small",
		"UPPER":         "uPPER",
		"UsersALL_CAPS": "usersALLCAPS",
		"snake_case":    "snakeCase",
		"camelCase":     "camelCase",
	}

	for in, out := range names {
		t.Run("should convert name "+in, func(t *testing.T) {
			r.Equal(t, out, PrivateGoName(in))
		})
	}
}

func TestDBName(t *testing.T) {
	names := map[string]string{
		"Same":          "same",
		"small":         "small",
		"UPPER":         "upper",
		"UsersALL_CAPS": "users_all_caps",
		"snake_case":    "snake_case",
		"camelCase":     "camel_case",
	}

	for in, out := range names {
		t.Run("should convert name "+in, func(t *testing.T) {
			r.Equal(t, out, DBName(in))
		})
	}
}

func TestJSONName(t *testing.T) {
	names := map[string]string{
		"Same":          "same",
		"small":         "small",
		"UPPER":         "upper",
		"UsersALL_CAPS": "usersAllCaps",
		"snake_case":    "snakeCase",
		"camelCase":     "camelCase",
	}

	for in, out := range names {
		t.Run("should convert name "+in, func(t *testing.T) {
			r.Equal(t, out, JSONName(in))
		})
	}
}
