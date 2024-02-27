package parser

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/ch3nnn/sql2pb/cmd/generation/tools/stringx"
)

var (

	// indent represents the indentation amount for fields. the style guide suggests
	// two spaces
	indent = "  "

	// gen protobuf field style
	fieldStyleToCamelWithStartLower = "sqlPb"
	fieldStyleToSnake               = "sql_pb"
)

type Message struct {
	Name    string
	Comment string
	Fields  []MessageField
	Style   string
}

// GenDefaultMessage gen default message
func (m *Message) GenDefaultMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	var curFields []MessageField
	var filedTag int
	for _, field := range m.Fields {
		if slices.Contains([]string{"version", "del_state", "delete_time"}, field.Name) {
			continue
		}
		filedTag++
		field.tag = filedTag
		field.Name = stringx.From(field.Name).ToCamelWithStartLower()
		if m.Style == fieldStyleToSnake {
			field.Name = stringx.From(field.Name).ToSnake()
		}

		if field.Comment == "" {
			field.Comment = field.Name
		}
		curFields = append(curFields, field)
	}
	m.Fields = curFields
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

func (m *Message) GenDefaultFilterMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = mOrginName + "Filter"
	var curFields []MessageField
	var filedTag int
	for _, field := range m.Fields {
		if slices.Contains([]string{"version", "del_state", "delete_time"}, field.Name) {
			continue
		}
		filedTag++
		field.tag = filedTag
		field.Name = stringx.From(field.Name).ToCamelWithStartLower()
		if m.Style == fieldStyleToSnake {
			field.Name = stringx.From(field.Name).ToSnake()
		}

		if field.Comment == "" {
			field.Comment = field.Name
		}
		// 可选
		field.Typ = "optional " + field.Typ

		curFields = append(curFields, field)
	}
	m.Fields = curFields
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

// GenRpcAddReqRespMessage gen add req message
func (m *Message) GenRpcAddReqRespMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	// req
	m.Name = "Add" + mOrginName + "Req"
	var curFields []MessageField
	var filedTag int
	for _, field := range m.Fields {
		if slices.Contains([]string{"id", "create_at", "create_time", "update_time", "update_at", "version", "del_state", "delete_time", "delete_at"}, field.Name) {
			continue
		}
		filedTag++
		field.tag = filedTag
		field.Name = stringx.From(field.Name).ToCamelWithStartLower()
		if m.Style == fieldStyleToSnake {
			field.Name = stringx.From(field.Name).ToSnake()
		}
		if field.Comment == "" {
			field.Comment = field.Name
		}
		curFields = append(curFields, field)
	}
	m.Fields = curFields
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	// resp
	m.Name = "Add" + mOrginName + "Resp"
	m.Fields = []MessageField{}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields

}

// GenRpcUpdateReqMessage gen add resp message
func (m *Message) GenRpcUpdateReqMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = "Update" + mOrginName + "Req"
	var curFields []MessageField
	var filedTag int
	for _, field := range m.Fields {
		if slices.Contains([]string{"create_time", "create_at", "update_time", "update_at", "version", "del_state", "delete_time", "delete_at"}, field.Name) {
			continue
		}
		filedTag++
		field.tag = filedTag
		field.Name = stringx.From(field.Name).ToCamelWithStartLower()
		if m.Style == fieldStyleToSnake {
			field.Name = stringx.From(field.Name).ToSnake()
		}
		if field.Comment == "" {
			field.Comment = field.Name
		}
		// 可选
		field.Typ = "optional " + field.Typ

		curFields = append(curFields, field)
	}
	m.Fields = curFields
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	// resp
	m.Name = "Update" + mOrginName + "Resp"
	m.Fields = []MessageField{}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

// GenRpcDelReqMessage gen add resp message
func (m *Message) GenRpcDelReqMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = "Del" + mOrginName + "Req"
	m.Fields = []MessageField{
		{Name: "id", Typ: "int64", tag: 1, Comment: "id"},
	}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	// resp
	m.Name = "Del" + mOrginName + "Resp"
	m.Fields = []MessageField{}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

// GenRpcGetByIdReqMessage gen add resp message
func (m *Message) GenRpcGetByIdReqMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = "Select" + mOrginName + "ByIdReq"
	m.Fields = []MessageField{
		{Name: "id", Typ: "int64", tag: 1, Comment: "id"},
	}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	// resp
	firstWord := strings.ToLower(string(m.Name[0]))
	m.Name = "Select" + mOrginName + "ByIdResp"

	name := stringx.From(firstWord + mOrginName[1:]).ToCamelWithStartLower()
	comment := stringx.From(firstWord + mOrginName[1:]).ToCamelWithStartLower()
	if m.Style == fieldStyleToSnake {
		name = stringx.From(firstWord + mOrginName[1:]).ToSnake()
		comment = stringx.From(firstWord + mOrginName[1:]).ToSnake()
	}
	m.Fields = []MessageField{
		{Typ: mOrginName, Name: name, tag: 1, Comment: comment},
	}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

// GenRpcSearchReqMessage gen add resp message
func (m *Message) GenRpcSearchReqMessage(buf *bytes.Buffer) {
	mOrginName := m.Name
	mOrginFields := m.Fields

	m.Name = "Select" + mOrginName + "ListReq"
	m.Fields = []MessageField{
		{Typ: "int64", Name: "page", tag: 1, Comment: "页码"},
		{Typ: "int64", Name: "page_size", tag: 2, Comment: "每页数量"},
		{Typ: "optional " + mOrginName + "Filter", Name: "filter", tag: 3, Comment: mOrginName + "Filter"},
	}

	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields

	// resp
	firstWord := strings.ToLower(string(m.Name[0]))
	m.Name = "Select" + mOrginName + "ListResp"

	// name := stringx.From(firstWord + mOrginName[1:]).ToCamelWithStartLower()
	comment := stringx.From(firstWord + mOrginName[1:]).ToCamelWithStartLower()
	if m.Style == fieldStyleToSnake {
		// name = stringx.From(firstWord + mOrginName[1:]).ToSnake()
		comment = stringx.From(firstWord + mOrginName[1:]).ToSnake()
	}

	m.Fields = []MessageField{
		{Typ: "int64", Name: "count", tag: 1, Comment: "总数"},
		{Typ: "int64", Name: "page_count", tag: 2, Comment: "页码总数"},
		{Typ: "repeated " + mOrginName, Name: "results", tag: 3, Comment: comment},
	}
	buf.WriteString(fmt.Sprintf("%s\n", m))

	// reset
	m.Name = mOrginName
	m.Fields = mOrginFields
}

// String returns a string representation of a Message.
func (m *Message) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("message %s {\n", m.Name))
	for _, f := range m.Fields {
		buf.WriteString(fmt.Sprintf("%s%s; //%s\n", indent, f, f.Comment))
	}
	buf.WriteString("}\n")

	return buf.String()
}

// AppendField appends a message field to a message. If the tag of the message field is in use, an error will be returned.
func (m *Message) AppendField(mf MessageField) error {
	for _, f := range m.Fields {
		if f.Tag() == mf.Tag() {
			return fmt.Errorf("tag `%d` is already in use by field `%s`", mf.Tag(), f.Name)
		}
	}

	m.Fields = append(m.Fields, mf)

	return nil
}

type MessageField struct {
	Typ     string
	Name    string
	tag     int
	Comment string
}

// NewMessageField creates a new message field.
func NewMessageField(typ, name string, tag int, comment string) MessageField {
	return MessageField{typ, name, tag, comment}
}

// Tag returns the unique numbered tag of the message field.
func (f MessageField) Tag() int {
	return f.tag
}

// String returns a string representation of a message field.
func (f MessageField) String() string {
	return fmt.Sprintf("%s %s = %d", f.Typ, f.Name, f.tag)
}
