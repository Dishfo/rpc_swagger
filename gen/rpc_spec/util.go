package rpc_spec

import (
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

	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, "-", "")
	return name
}

func NeedOmit(field reflect.StructField) bool {
	if field.Name[0] < 'z' && field.Name[0] > 'a' {
		return true
	}

	val := field.Tag.Get("json")
	return val == "-"
}
