package dagger

import (
	"reflect"
	"strings"
	"unicode"
)

// Returns true if the given method is exported, i.e. it starts with an
// uppercase character.
func isMethodExported(m reflect.Method) bool {
	return unicode.IsUpper(rune(m.Name[0]))
}

// Returns true if the given method is a provider method, i.e. it starts with
// the string "Provide".
func isProviderMethod(m reflect.Method) bool {
	return strings.HasPrefix(m.Name, "Provide")
}
