package banister

import (
	"fmt"
	"strings"
)

func BuildTableSQL(b Backend, m Model) string {
	definitionSQL := make([]string, 0)
	for _, f := range m.Fields() {
		columnSQL := BuildColumnSQL(b, f, true)
		if columnSQL == "" {
			continue
		}

		definitionSQL = append(definitionSQL, columnSQL)
	}

	constraints := b.ConstraintSQL(m)
	if constraints != "" {
		definitionSQL = append(definitionSQL, b.ConstraintSQL(m))
	}

	tableSQL := b.MigrationOperations().CreateTable
	tableSQL = strings.ReplaceAll(tableSQL, "__TABLE__", m.Settings().DBTable)
	tableSQL = strings.ReplaceAll(tableSQL, "__DEFINITION__", strings.Join(definitionSQL, ", "))

	return tableSQL + ";"
}

func BuildColumnSQL(b Backend, f Field, includeDefault bool) string {
	sqlType, typeExists := b.DataTypes()[f.Type()]
	if !typeExists {
		panic("Type not supported: " + f.Type())
	}

	settings := f.Settings()
	columnSQL := fmt.Sprintf("%s %s", settings.DBColumn, sqlType)

	if includeDefault && settings.Default.IsValid() {
		columnSQL += " DEFAULT " + b.QuoteValue(settings.Default.Value)
	}

	// Set NULL constraint
	if settings.PrimaryKey {
		// Primary keys are always NOT NULL
	} else if settings.Null {
		// Could set this but it's the default anyway
		// columnSQL += " NULL"
	} else {
		columnSQL += " NOT NULL"
	}

	// Set uniqueness
	// NOTE: if it's a PK then UNIQUE is redundant
	if settings.PrimaryKey {
		columnSQL += " PRIMARY KEY"
	} else if settings.Unique {
		columnSQL += " UNIQUE"
	}

	if suffix, ok := b.DataTypeSuffixes()[f.Type()]; ok {
		columnSQL += " " + suffix
	}

	if strings.Contains(columnSQL, "__MAX_LENGTH__") {
		if f.Settings().MaxLength == nil {
			panic("Field " + settings.Name + " requires max length to be set")
		}
		columnSQL = strings.ReplaceAll(columnSQL, "__MAX_LENGTH__", fmt.Sprintf("%d", *settings.MaxLength))
	}

	if strings.Contains(columnSQL, "__SIZE__") {
		if settings.Size != nil {
			columnSQL = strings.ReplaceAll(columnSQL, "__SIZE__", fmt.Sprintf("%d", settings.Size))
		} else {
			columnSQL = strings.ReplaceAll(columnSQL, "__SIZE__", "")
		}
	}

	return columnSQL
}

func BuildConstraintSQL(b Backend, m Model) string {
	definitionSQL := make([]string, 0)

	for _, f := range m.Fields() {
		if f.Settings().Rel == nil {
			continue
		}

		v := b.MigrationOperations().CreateFK
		v = strings.ReplaceAll(v, "__COLUMN__", f.Settings().DBColumn)
		v = strings.ReplaceAll(v, "__TO_TABLE__", f.Settings().Rel.To.Settings().DBTable)
		v = strings.ReplaceAll(v, "__TO_COLUMN__", f.Settings().Rel.ToField.Settings().DBColumn)
		if f.Settings().Rel.OnDelete != "" {
			v = strings.ReplaceAll(v, "__ON_DELETE__", " ON DELETE "+f.Settings().Rel.OnDelete)
		} else {
			v = strings.ReplaceAll(v, "__ON_DELETE__", "")
		}
		definitionSQL = append(definitionSQL, v)
	}

	return strings.Join(definitionSQL, ", ")
}
