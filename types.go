package coi

import "reflect"

func getRawType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Ptr:
		return getRawType(t.Elem())
	case reflect.Array:
		return getRawType(t.Elem())
	case reflect.Chan:
		return getRawType(t.Elem())
	case reflect.Map:
		return getRawType(t.Elem())
	case reflect.Slice:
		return getRawType(t.Elem())
	default:
		return t
	}
}
