package banister

type Backend interface {
	DriverName() string
	DisplayName() string
	MigrationOperations() DBOperations
	Lookups() map[QueryOperator]string
	DataTypes() map[FieldType]string
	DataTypeSuffixes() map[FieldType]string
	ReturnInsertColumnsSQL(m Model) string
	ColumnSQL(f Field, includeDefault bool) string
	TableSQL(m Model) string
	ConstraintSQL(m Model) string
	QuoteName(name string) string
	QuoteValue(value interface{}) string
}

type DBOperations struct {
	CreateTable  string
	CreateColumn string
	CreateIndex  string
	CreateFK     string
}

var backends = map[string]Backend{}

func RegisterBackend(b Backend) {
	backends[b.DriverName()] = b
}

func GetBackend(driver string) Backend {
	if b, ok := backends[driver]; ok {
		return b
	}

	panic("backend not found for driver: " + driver + "\n  Did you forget to import it?")
}
