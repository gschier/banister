package banister_test

import (
	"github.com/dave/jennifer/jen"
	. "github.com/gschier/banister"
	"github.com/gschier/banister/testutil"
	assert "github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestModelGenerator_GoSource(t *testing.T) {
	file := jen.NewFile("dummy")
	NewModelGenerator(file, testutil.TestUserModel()).Generate()
	assert.Equal(t, strings.TrimSpace(`
package dummy

import (
	"encoding/json"
	"fmt"
	"time"
)

// User is a database model which represents a single row from the
// users database table
type User struct {
	ID       int64
	Age      *int64
	Name     string
	Username string
	Created  time.Time
}

// PrintJSON prints out a JSON string of the model for debugging
func (model *User) PrintJSON() {
	b, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T: %s", model, b)
}
`), strings.TrimSpace(file.GoString()))
}
