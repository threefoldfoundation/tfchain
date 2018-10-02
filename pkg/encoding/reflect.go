package encoding

import (
	"reflect"
	"unicode"
)

func isFieldHidden(val reflect.Value, index int) bool {
	field := val.Type().Field(index)
	if field.Anonymous {
		return true
	}
	for _, r := range field.Name {
		return r == '_' || unicode.IsLower(r)
	}
	return true
}
