package banister_test

import (
	"github.com/dave/jennifer/jen"
	. "github.com/gschier/banister"
	. "github.com/gschier/banister/testutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFilterGenerator_GoSource(t *testing.T) {
	file := jen.NewFile("dummy")
	m := TestUserModel()
	NewFilterGenerator(file, m.Fields()[1], m).Generate()
	assert.Equal(t, strings.TrimSpace(`
package dummy

import squirrel "github.com/Masterminds/squirrel"

type UserAgeFilter struct {
	filters []squirrel.Sqlizer
}

func (filter *UserAgeFilter) Eq(v int64) userQuerysetFilterArg {
	return userQuerysetFilterArg{filter: &squirrel.Eq{"\"users\".\"age\"": v}}
}
func (filter *UserAgeFilter) Lt(v int64) userQuerysetFilterArg {
	return userQuerysetFilterArg{filter: &squirrel.Lt{"\"users\".\"age\"": v}}
}
func (filter *UserAgeFilter) Lte(v int64) userQuerysetFilterArg {
	return userQuerysetFilterArg{filter: &squirrel.LtOrEq{"\"users\".\"age\"": v}}
}
func (filter *UserAgeFilter) Gt(v int64) userQuerysetFilterArg {
	return userQuerysetFilterArg{filter: &squirrel.Gt{"\"users\".\"age\"": v}}
}
func (filter *UserAgeFilter) Gte(v int64) userQuerysetFilterArg {
	return userQuerysetFilterArg{filter: &squirrel.GtOrEq{"\"users\".\"age\"": v}}
}
func (filter *UserAgeFilter) Null() userQuerysetFilterArg {
	return userQuerysetFilterArg{filter: &squirrel.Eq{"\"users\".\"age\"": nil}}
}
`), strings.TrimSpace(file.GoString()))
}
