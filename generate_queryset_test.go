package banister_test

import (
	"github.com/dave/jennifer/jen"
	. "github.com/gschier/banister"
	. "github.com/gschier/banister/testutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestQuerysetGenerator_AddStarSelectMethod(t *testing.T) {
	file := jen.NewFile("dummy")
	m := TestUserModel()
	NewQuerysetGenerator(file, m).AddStarSelectMethod()
	assert.Equal(t, strings.TrimSpace(`
package dummy

import squirrel "github.com/Masterminds/squirrel"

func (qs *UserQueryset) starSelect() squirrel.SelectBuilder {
	query := squirrel.Select(
		"\"users\".\"id\"",
		"\"users\".\"age\"",
		"\"users\".\"name\"",
		"\"users\".\"username\"",
		"\"users\".\"created\"").From("users")

	joinCheck := make(map[string]bool)

	// Assign filters and join if necessary
	for _, w := range qs.filter {
		query = query.Where(w.filter)
		for _, j := range w.joins {
			if _, ok := joinCheck[j]; ok {
				continue
			}
			joinCheck[j] = true
			query = query.Join(j)
		}
	}

	// Apply limit if set
	if qs.limit > 0 {
		query = query.Limit(f.limit)
	}

	// Apply offset if set
	if qs.offset > 0 {
		query = query.Offset(f.offset)
	}

	// Apply default order if none specified
	if len(qs.orderBy) == 0 {
		// TODO: Add default order-by
	}

	// Apply user-specified order
	for _, s := range qs.orderBy {
		query = query.OrderBy(s.field + " " + s.order)
	}

	return query
}
`), strings.TrimSpace(file.GoString()))
}

func TestQuerysetGenerator_AddDeleteMethod(t *testing.T) {
	file := jen.NewFile("dummy")
	m := TestUserModel()
	NewQuerysetGenerator(file, m).AddDeleteMethod()
	assert.Equal(t, strings.TrimSpace(`
package dummy

import squirrel "github.com/Masterminds/squirrel"

func (qs *UserQueryset) Delete(m *User) error {
	query := squirrel.Delete("users")

	for _, w := range qs.filters {
		query = query.Where(w.filter)
	}

	q, args := toSQL(query)
	_, err := f.mgr.db.Exec(q, args...)
	return err
}
`), strings.TrimSpace(file.GoString()))
}
