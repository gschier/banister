package banister_test

import (
	"github.com/dave/jennifer/jen"
	. "github.com/gschier/banister"
	. "github.com/gschier/banister/testutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestQuerysetGenerator_Generate(t *testing.T) {
	file := jen.NewFile("dummy")
	m := TestUserModel()
	NewQuerysetGenerator(file, m).Generate()
	assert.Equal(t, strings.TrimSpace(`
package dummy

import squirrel "github.com/Masterminds/squirrel"

type UserQueryset struct {
	filter  []userQuerysetFilterArg
	orderBy []userQuerysetOrderByArg
	limit   uint64
	offset  uint64
}

func NewUserQueryset() *UserQueryset {
	return &UserQueryset{
		filter:  make([]userQuerysetFilterArg, 0),
		limit:   0,
		offset:  0,
		orderBy: make([]userQuerysetOrderByArg, 0),
	}
}
func (qs *UserQueryset) Filter(filter ...userQuerysetFilterArg) *UserQueryset {
	qs.filter = append(qs.filter, filter...)
	return qs
}
func (qs *UserQueryset) Order(orderBy ...userQuerysetOrderByArg) *UserQueryset {
	qs.orderBy = append(qs.orderBy, orderBy...)
	return qs
}
func (qs *UserQueryset) Limit(limit uint64) *UserQueryset {
	qs.limit = limit
	return qs
}
func (qs *UserQueryset) Offset(offset uint64) *UserQueryset {
	qs.offset = offset
	return qs
}

type userQuerysetFilterArg struct {
	filter squirrel.Sqlizer
	joins  []string
}
type userQuerysetOrderByArg struct {
	field string
	order string
	join  string
}
type userQuerysetSetterArg struct {
	field string
	value interface{}
}

var WhereUser = struct {
	ID       UserIDFilter
	Age      UserAgeFilter
	Name     UserNameFilter
	Username UserUsernameFilter
	Created  UserCreatedFilter
}{
	Age:      &UserAgeFilter{},
	Created:  &UserCreatedFilter{},
	ID:       &UserIDFilter{},
	Name:     &UserNameFilter{},
	Username: &UserUsernameFilter{},
}
`), strings.TrimSpace(file.GoString()))
}
