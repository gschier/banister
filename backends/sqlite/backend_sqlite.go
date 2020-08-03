package sqlite

import (
	"fmt"
	"github.com/gschier/banister"
	"reflect"
	"strings"
)

func init() {
	banister.RegisterBackend(&Backend{})
}

type Backend struct {
}

func (b *Backend) DriverName() string {
	return "sqlite3"
}

func (b *Backend) DisplayName() string {
	return "SQLite"
}

func (b *Backend) DataTypes() map[banister.FieldType]string {
	return map[banister.FieldType]string{
		banister.Auto:      "INTEGER",
		banister.Boolean:   "BOOLEAN",
		banister.Char:      "VARCHAR(__MAX_LENGTH__)",
		banister.DateTime:  "DATETIME",
		banister.Duration:  "INTERVAL",
		banister.Float:     "REAL",
		banister.Integer:   "INTEGER",
		banister.Text:      "TEXT",
		banister.TextArray: "TEXT",
	}
}

func (b *Backend) DataTypeSuffixes() map[banister.FieldType]string {
	return map[banister.FieldType]string{
		banister.Auto: "AUTOINCREMENT",
	}
}

func (b *Backend) ReturnInsertColumnsSQL(m banister.Model) string {
	return ""
}

func (b *Backend) Lookups() map[banister.QueryOperator]string {
	return map[banister.QueryOperator]string{
		banister.Exact:       "= ?",
		banister.NotExact:    "!= ?",
		banister.IExact:      `LIKE ? ESCAPE '\'`,
		banister.Contains:    `LIKE '%' || ? || '%' ESCAPE '\'`,
		banister.IContains:   `ILIKE '%' || ? || '%' ESCAPE '\'`,
		banister.Regex:       "REGEXP ?",
		banister.IRegex:      "REGEXP '(??i)' || ?",
		banister.Gt:          "> ?",
		banister.Gte:         ">= ?",
		banister.Lt:          "< ?",
		banister.Lte:         "<= ?",
		banister.StartsWith:  `LIKE ? || '%' ESCAPE '\'`,
		banister.IStartsWith: `ILIKE ? || '%' ESCAPE '\'`,
		banister.EndsWith:    `LIKE '%' || ? ESCAPE '\'`,
		banister.IEndsWith:   `ILIKE '%' || ? ESCAPE '\'`,
	}
}

func (b *Backend) ColumnSQL(f banister.Field, includeDefault bool) string {
	return banister.BuildColumnSQL(b, f, includeDefault)
}

func (b *Backend) TableSQL(m banister.Model) string {
	return banister.BuildTableSQL(b, m)
}

func (b *Backend) MigrationOperations() banister.DBOperations {
	return banister.DBOperations{
		CreateTable:  "CREATE TABLE __TABLE__ ( __DEFINITION__ )",
		CreateColumn: "ALTER TABLE __TABLE__ ADD COLUMN __COLUMN__ __DEFINITION__",
		CreateIndex:  "CREATE INDEX __NAME__ ON __TABLE__ (__COLUMN__)__INCLUDE____EXTRA____CONDITION__",
		CreateFK:     "FOREIGN KEY (__COLUMN__) REFERENCES __TO_TABLE__ (__TO_COLUMN__)__ON_DELETE__",
	}
}

func (b *Backend) ConstraintSQL(m banister.Model) string {
	return banister.BuildConstraintSQL(b, m)
}

func (b *Backend) QuoteName(name string) string {
	// Quoting once is enough
	if strings.HasPrefix(name, `"`) && strings.HasSuffix(name, `"`) {
		return name
	}
	return fmt.Sprintf(`"%s"`, name)
}

func (b *Backend) QuoteValue(value interface{}) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case int, int64, int32:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case bool:
		if v {
			return "1"
		} else {
			return "0"
		}
	default:
		panic("Cannot quote unsupported Go type: " + fmt.Sprintf("%s", reflect.TypeOf(v)))
	}
}
