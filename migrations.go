package banister

import (
	"fmt"
	"strings"
)

type Migration struct {
	Models []Model
}

func (mig *Migration) model(name string) Model {
	for _, m := range mig.Models {
		if strings.ToLower(m.Settings().Name) == strings.ToLower(name) {
			return m
		}
	}

	panic("failed to find model by name: " + name)
}

func (mig *Migration) field(m Model, name string) Field {
	for _, f := range m.Fields() {
		if strings.ToLower(f.Settings().Name) == strings.ToLower(name) {
			return f
		}
	}

	panic("failed to find field in " + m.Settings().Name + " model by name: " + name)
}

func (mig *Migration) RenameField(model, old, new string) string {
	m := mig.model(model)
	f := mig.field(m, old)

	return fmt.Sprintf("RENAME COLUMN %s TO %s;", f.Settings().DBColumn, new)
}
