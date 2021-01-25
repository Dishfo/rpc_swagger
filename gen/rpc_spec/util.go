package rpc_spec

import (
	"encoding/json"
	"fmt"
	"reflect"
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
