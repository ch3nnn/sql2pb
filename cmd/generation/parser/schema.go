package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/chuckpreslar/inflect"
	"github.com/serenize/snaker"
	"github.com/sirupsen/logrus"
)

type Schema struct {
	Syntax      string
	ServiceName string
	GoPackage   string
	Package     string
	Imports     sort.StringSlice
	Messages    []*Message
	Enums       []*Enum
}

func NewSchema(syntax string, serviceName string, goPackage string, Package string) *Schema {
	return &Schema{Syntax: syntax, ServiceName: serviceName, GoPackage: goPackage, Package: Package}
}

// TypesFromColumns creates the appropriate schema properties from a collection of column types.
func (s *Schema) TypesFromColumns(cols []Column, ignoreTables, ignoreColumns []string, fieldStyle string) error {
	messageMap := map[string]*Message{}
	ignoreMap := map[string]bool{}
	ignoreColumnMap := map[string]bool{}
	for _, ig := range ignoreTables {
		ignoreMap[ig] = true
	}
	for _, ic := range ignoreColumns {
		ignoreColumnMap[ic] = true
	}

	for _, c := range cols {
		if _, ok := ignoreMap[c.TableName]; ok {
			continue
		}
		if _, ok := ignoreColumnMap[c.ColumnName]; ok {
			continue
		}

		messageName := snaker.SnakeToCamel(c.TableName)

		msg, ok := messageMap[messageName]
		if !ok {
			messageMap[messageName] = &Message{Name: messageName, Comment: c.TableComment, Style: fieldStyle}
			msg = messageMap[messageName]
		}

		err := s.parseColumn(msg, c)
		if nil != err {
			return err
		}
	}

	for _, v := range messageMap {
		s.Messages = append(s.Messages, v)
	}

	return nil
}

// parseColumn parses a column and inserts the relevant fields in the Message. If an enumerated type is encountered, an Enum will
// be added to the Schema. Returns an error if an incompatible protobuf data type cannot be found for the database column type.
func (s *Schema) parseColumn(msg *Message, col Column) error {
	typ := strings.ToLower(col.DataType)
	var fieldType string

	switch typ {
	case "char", "varchar", "text", "longtext", "mediumtext", "tinytext":
		fieldType = "string"
	case "enum", "set":
		// Parse c.ColumnType to get the enum list
		enumList := regexp.MustCompile(`[enum|set]\((.+?)\)`).FindStringSubmatch(col.ColumnType)
		enums := strings.FieldsFunc(enumList[1], func(c rune) bool {
			cs := string(c)
			return "," == cs || "'" == cs
		})

		enumName := inflect.Singularize(snaker.SnakeToCamel(col.TableName)) + snaker.SnakeToCamel(col.ColumnName)
		enum, err := newEnumFromStrings(enumName, col.ColumnComment, enums)
		if nil != err {
			return err
		}

		s.Enums = append(s.Enums, enum)

		fieldType = enumName
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		fieldType = "bytes"
	case "date", "time", "datetime", "timestamp", "timestamptz":
		// s.AppendImport("google/protobuf/timestamp.proto")
		fieldType = "int64"
	case "bool", "bit":
		fieldType = "bool"
	case "tinyint", "smallint", "int", "mediumint", "bigint", "int2", "int4", "int8":
		if col.ColumnType == "tinyint(1)" {
			fieldType = "bool"
			break
		}
		fieldType = "int64"
	case "float", "decimal", "double":
		fieldType = "double"
	case "json":
		fieldType = "string"
	}

	if "" == fieldType {
		fieldType = "string"
		logrus.Warning(fmt.Errorf("no compatible protobuf type found for `%s`. column: `%s`.`%s`. default set column 'string'", col.DataType, col.TableName, col.ColumnName).Error())
	}

	field := NewMessageField(fieldType, col.ColumnName, len(msg.Fields)+1, col.ColumnComment)

	err := msg.AppendField(field)
	if nil != err {
		return err
	}

	return nil
}

// String returns a string representation of a Schema.
func (s *Schema) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("syntax = \"%s\";\n", s.Syntax))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("option go_package =\"%s\";\n", s.GoPackage))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("package %s;\n", s.Package))

	buf.WriteString("\n")
	buf.WriteString("// ------------------------------------ \n")
	buf.WriteString("// Messages\n")
	buf.WriteString("// ------------------------------------ \n\n")

	for _, m := range s.Messages {
		buf.WriteString("//--------------------------------" + m.Comment + "--------------------------------")
		buf.WriteString("\n\n")
		m.GenDefaultMessage(buf)
		m.GenDefaultFilterMessage(buf)
		m.GenRpcAddReqRespMessage(buf)
		m.GenRpcUpdateReqMessage(buf)
		m.GenRpcDelReqMessage(buf)
		m.GenRpcGetByIdReqMessage(buf)
		m.GenRpcSearchReqMessage(buf)
	}

	buf.WriteString("\n")

	if len(s.Enums) > 0 {
		buf.WriteString("// ------------------------------------ \n")
		buf.WriteString("// Enums\n")
		buf.WriteString("// ------------------------------------ \n\n")

		for _, e := range s.Enums {
			buf.WriteString(fmt.Sprintf("%s\n", e))
		}
	}

	buf.WriteString("\n")
	buf.WriteString("// ------------------------------------ \n")
	buf.WriteString("// Rpc Func\n")
	buf.WriteString("// ------------------------------------ \n\n")

	funcTpl := "service " + s.ServiceName + "{ \n\n"
	for _, m := range s.Messages {
		funcTpl += "\t //-----------------------" + m.Comment + "----------------------- \n"
		funcTpl += "\n\t // 创建" + m.Comment + "\n"
		funcTpl += "\t rpc Insert" + m.Name + "(Add" + m.Name + "Req) returns (Add" + m.Name + "Resp); \n"
		funcTpl += "\n\t // 更新" + m.Comment + "\n"
		funcTpl += "\t rpc Update" + m.Name + "(Update" + m.Name + "Req) returns (Update" + m.Name + "Resp); \n"
		funcTpl += "\n\t // 根据 " + m.Comment + " id 删除\n"
		funcTpl += "\t rpc Delete" + m.Name + "(Del" + m.Name + "Req) returns (Del" + m.Name + "Resp); \n"
		funcTpl += "\n\t // 根据 " + m.Comment + " id 获取详情\n"
		funcTpl += "\t rpc Select" + m.Name + "ById(Select" + m.Name + "ByIdReq) returns (Select" + m.Name + "ByIdResp); \n"
		funcTpl += "\n\t // " + m.Comment + " 列表\n"
		funcTpl += "\t rpc Select" + m.Name + "List(Select" + m.Name + "ListReq) returns (Select" + m.Name + "ListResp); \n"
	}
	funcTpl = funcTpl + "\n}"
	buf.WriteString(funcTpl)

	return buf.String()
}
