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

	tableSQL := b.Operations().CreateTable
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
		columnSQL += " DEFAULT " + b.QuoteValue(settings.Default)
	}

	if settings.Null {
		columnSQL += " NULL"
	} else {
		columnSQL += " NOT NULL"
	}

	if settings.PrimaryKey {
		columnSQL += " PRIMARY KEY"
	} else if settings.Unique {
		columnSQL += " UNIQUE"
	}

	if suffix, ok := b.DataTypeSuffixes()[f.Type()]; ok {
		columnSQL += " " + suffix
	}

	if strings.Contains(columnSQL, "__MAX_LENGTH__") {
		if f.Settings().MaxLength != nil {
			panic("Field " + settings.Name + "requires max length to be set")
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
	// TODO: Add support for model-wide constraints
	return strings.Join(definitionSQL, ", ")
}
