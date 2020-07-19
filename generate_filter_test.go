package banister_test

import (
	"github.com/dave/jennifer/jen"
	. "github.com/gschier/banister"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFilterGenerator_GoSource(t *testing.T) {
	file := jen.NewFile("dummy")
	NewFilterGenerator(file, NewIntegerField("age").Null().Build()).Generate()
	assert.Equal(t, strings.TrimSpace(`
package dummy

import squirrel "github.com/Masterminds/squirrel"

type AgeFilter struct {
	filters []squirrel.Sqlizer
}

// Eq does a thing
func (filter *AgeFilter) Eq(v string) *AgeFilter {
	filter.filters = append(filter, &squirrel.Eq{"age": "v"})
	return filter
}

// Gt does a thing
func (filter *AgeFilter) Gt(v string) *AgeFilter {
	filter.filters = append(filter, &squirrel.Gt{"age": "v"})
	return filter
}

// Gte does a thing
func (filter *AgeFilter) Gte(v string) *AgeFilter {
	filter.filters = append(filter, &squirrel.GtOrEq{"age": "v"})
	return filter
}

// Lt does a thing
func (filter *AgeFilter) Lt(v string) *AgeFilter {
	filter.filters = append(filter, &squirrel.Lt{"age": "v"})
	return filter
}

// Lte does a thing
func (filter *AgeFilter) Lte(v string) *AgeFilter {
	filter.filters = append(filter, &squirrel.LtOrEq{"age": "v"})
	return filter
}

// Null does a thing
func (filter *AgeFilter) Null() *AgeFilter {
	filter.filters = append(filter, &squirrel.Eq{"age": nil})
	return filter
}
`), strings.TrimSpace(file.GoString()))
}
