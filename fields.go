package banister

type FieldSettings struct {
	Name         string
	Hidden       bool
	VerboseName  string
	PrimaryKey   bool
	Unique       bool
	Null         bool
	DBIndex      bool
	Rel          *Rel
	Default      DefaultValue
	Choices      *map[string]interface{}
	DBColumn     string
	DBTablespace string
	Validators   []FieldValidator

	// MaxLength specifies the length of a char field
	MaxLength *uint

	// Size specifies the size of an array field
	Size *uint
}

func (s *FieldSettings) Fix() {
	s.Name = PublicGoName(s.Name)

	if s.DBColumn == "" {
		s.DBColumn = DBName(s.Name)
	}
}

func (s FieldSettings) Names(modelName string) GeneratedFieldNames {
	return NamesForField(modelName, s.Name)
}

func (s FieldSettings) GenFilterStructName(modelName string) string {
	return PrivateGoName(modelName + s.Name + "Filter")
}

func (s FieldSettings) GenOrderStructName(modelName string) string {
	return PrivateGoName(modelName + s.Name + "OrderBy")
}

func (s FieldSettings) GenSetStructName(modelName string) string {
	return PrivateGoName(modelName + s.Name + "Set")
}

type DefaultValue struct {
	Value interface{}
	valid bool
}

func NewDefaultValue(v interface{}) DefaultValue {
	return DefaultValue{valid: true, Value: v}
}

func (d DefaultValue) IsValid() bool {
	return d.valid
}

type Rel struct {
	// To specifies the name of the model that the relationship points to
	To string

	// ToField specifies the field of the related model that this relation
	// links to.
	//
	// This is usually the primary key, which is usually ID
	ToField string

	// RelatedName specifies the reverse name that links back to this
	// field's model
	//
	// For example, if a Post has an FK to User, the related name "Posts"
	// would specify how the relation would be accessed from User.
	RelatedName string

	// RelatedQueryName specifies the reverse name when used in queries.
	//
	// For example, a related name of Posts might have a related query name
	// of Post
	RelatedQueryName string

	// OnDelete specifies what should happen when a model is deleted
	// that has links to itself via foreign keys
	// TODO: Make this an enum (eg. CASCADE | DO NOTHING | SET DEFAULT)
	OnDelete string
}

type FieldType string

const (
	Auto      FieldType = "Auto"
	Char      FieldType = "Char"
	DateTime  FieldType = "DateTime"
	Duration  FieldType = "Duration"
	Integer   FieldType = "Integer"
	Text      FieldType = "Text"
	TextArray FieldType = "TextArray"
	Float     FieldType = "Float"
	Boolean   FieldType = "Boolean"
)

type FieldValidator func(v interface{}) error

type Field interface {
	Type() FieldType
	RelType() FieldType
	Settings() FieldSettings
	EmptyDefault() interface{}
	Operations() []Operation
}
