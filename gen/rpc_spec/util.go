package rpc_spec

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func indirectType(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Ptr:
		return indirectType(typ.Elem())
	default:
		return typ
	}
}

func TypeName(typ reflect.Type) string {
	name := typ.String()

	//name = strings.ReplaceAll(name, ".", "")
	//name = strings.ReplaceAll(name, "_", "")
	//name = strings.ReplaceAll(name, "-", "")
	return name
}

func NeedOmit(field reflect.StructField) bool {
	if field.Name[0] < 'z' && field.Name[0] > 'a' {
		return true
	}

	val := field.Tag.Get("json")
	return val == "-"
}

func isStringer(typ reflect.Type) bool {
	var stringer *fmt.Stringer = nil

	return typ.Implements(reflect.TypeOf(stringer).Elem())
}

func isMarshal(typ reflect.Type) bool {
	var stringer *json.Marshaler = nil

	return typ.Implements(reflect.TypeOf(stringer).Elem())
}

func isUnMarshal(typ reflect.Type) bool {
	var stringer *json.Unmarshaler = nil

	return typ.Implements(reflect.TypeOf(stringer).Elem())
}

func FieldName(fieldType reflect.StructField) string {
	val := fieldType.Tag.Get("json")
	if val == "-" {
		return fieldType.Name
	}
	if len(val) > 0 {
		tagItems := strings.Split(val, ",")
		if tagItems[0] != "-" {
			return tagItems[0]
		}
	}

	return fieldType.Name
}

func PackingDescription(comments []string) string {
	var builder strings.Builder
	for _, comment := range comments {
		tmp := strings.TrimLeft(comment, "//")
		tmp = strings.TrimRight(tmp, "\n")
		tmp = strings.TrimLeft(tmp, "/*")
		tmp = strings.TrimRight(tmp, "*/")
		tmp = strings.TrimSpace(tmp)
		builder.WriteString(tmp)
		builder.WriteRune(' ')
	}

	res := builder.String()
	res = strings.TrimLeft(res, "/*")
	res = strings.TrimRight(res, "*/")
	return strings.ReplaceAll(strings.TrimSpace(res), "\n", " ")
}
