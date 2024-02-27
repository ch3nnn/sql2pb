package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

//-----------------------Enum--------------------------------

// Enum represents a protocol buffer enumerated type.
type Enum struct {
	Name    string
	Comment string
	Fields  []EnumField
}

// String returns a string representation of an Enum.
func (e *Enum) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("// %s \n", e.Comment))
	buf.WriteString(fmt.Sprintf("enum %s {\n", e.Name))

	for _, f := range e.Fields {
		buf.WriteString(fmt.Sprintf("%s%s;\n", indent, f))
	}

	buf.WriteString("}\n")

	return buf.String()
}

// AppendField appends an EnumField to an Enum.
func (e *Enum) AppendField(ef EnumField) error {
	for _, f := range e.Fields {
		if f.Tag() == ef.Tag() {
			return fmt.Errorf("tag `%d` is already in use by field `%s`", ef.Tag(), f.Name())
		}
	}

	e.Fields = append(e.Fields, ef)

	return nil
}

//-----------------------EnumField--------------------------------

// EnumField represents a field in an enumerated type.
type EnumField struct {
	name string
	tag  int
}

// NewEnumField constructs an EnumField type.
func NewEnumField(name string, tag int) EnumField {
	name = strings.ToUpper(name)

	re := regexp.MustCompile(`(\W+)`)
	name = re.ReplaceAllString(name, "_")

	return EnumField{name, tag}
}

// String returns a string representation of an Enum.
func (ef EnumField) String() string {
	return fmt.Sprintf("%s = %d", ef.name, ef.tag)
}

// Name returns the name of the enum field.
func (ef EnumField) Name() string {
	return ef.name
}

// Tag returns the identifier tag of the enum field.
func (ef EnumField) Tag() int {
	return ef.tag
}

// newEnumFromStrings creates an enum from a name and a slice of strings that represent the names of each field.
func newEnumFromStrings(name, comment string, ss []string) (*Enum, error) {
	enum := &Enum{}
	enum.Name = name
	enum.Comment = comment

	for i, s := range ss {
		err := enum.AppendField(NewEnumField(s, i))
		if nil != err {
			return nil, err
		}
	}

	return enum, nil
}
