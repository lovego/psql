package bsql

import (
	"reflect"
	"regexp"
	"strings"
)

func Column2Field(column string) string {
	var parts []string
	for _, part := range strings.Split(column, "_") {
		parts = append(parts, strings.Title(part))
	}
	return strings.Join(parts, "")
}

var fieldToColumnRe1 = regexp.MustCompile("(.)([A-Z][a-z]+)")
var fieldToColumnRe2 = regexp.MustCompile("([a-z0-9])([A-Z])")

func Field2Column(field string) string {
	field = fieldToColumnRe1.ReplaceAllString(field, "$1_$2")
	field = fieldToColumnRe2.ReplaceAllString(field, "$1_$2")
	return strings.ToLower(field)
}

func Columns2Fields(columns []string) (result []string) {
	for _, column := range columns {
		result = append(result, Column2Field(column))
	}
	return result
}

func Fields2Columns(fields []string) (result []string) {
	for _, field := range fields {
		result = append(result, Field2Column(field))
	}
	return result
}

func FieldsFromStruct(v interface{}, exclude []string) (result []string) {
	LoopStructFields(reflect.ValueOf(v).Type(), func(name string) {
		if NotIn(name, exclude) {
			result = append(result, name)
		}
	})
	return
}

func LoopStructFields(typ reflect.Type, fn func(name string)) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			LoopStructFields(field.Type, fn)
		} else {
			fn(field.Name)
		}
	}
}

func NotIn(target string, slice []string) bool {
	for _, elem := range slice {
		if elem == target {
			return false
		}
	}
	return true
}
