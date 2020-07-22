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
	JSONName     string
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

	if s.JSONName == "" {
		s.JSONName = JSONName(s.Name)
	}
}

func (s FieldSettings) Names(model Model) GeneratedFieldNames {
	return NamesForField(model.Settings(), s)
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
	Settings() FieldSettings

	// Type specifies the type of field, which indirectly maps to database
	// column.
	Type() FieldType

	// RelType specifies the field type for FK's that need to reference this
	// field.
	//
	// For example, FK's that reference an Auto field will be stored as
	// Integer
	RelType() FieldType

	// EmptyDefault defines the default Go value to be used for non-pointer
	// Go variables.
	//
	// For example, instantiating a model struct.
	EmptyDefault() interface{}

	// Operations specifies a list of queryset operations that can be
	// performed on the field
	//
	// For example
	//   - Integer will have >=, >, ==, etc
	//   - Boolean will have ==, etc
	Operations() []Operation
}
