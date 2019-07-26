package autojson

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()
var headerType = reflect.TypeOf((*HeaderProvider)(nil)).Elem()

type signature struct {
	InParameters bool
	Status       bool
	Error        bool
}

func writeType(sb *strings.Builder, t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		sb.WriteString("*")
		t = t.Elem()
	}

	pkgpath := t.PkgPath()
	if pkgpath != "" {
		sb.WriteString(pkgpath[strings.LastIndex(pkgpath, "/")+1:])
		sb.WriteString(".")
	}

	sb.WriteString(t.Name())
}

func dumpSignature(t reflect.Type) string {
	sb := strings.Builder{}
	sb.WriteString("func(")

	for i := 0; i < t.NumIn(); i++ {
		writeType(&sb, t.In(i))

		if i < t.NumIn()-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(")")

	if t.NumOut() > 0 {
		sb.WriteString(" (")

		for i := 0; i < t.NumOut(); i++ {
			writeType(&sb, t.Out(i))

			if i < t.NumOut()-1 {
				sb.WriteString(", ")
			}
		}

		sb.WriteString(")")
	}

	return sb.String()
}

func validateSignature(h interface{}) (signature, error) {
	result := signature{}
	t := reflect.TypeOf(h)

	if t.Kind() != reflect.Func {
		return result, fmt.Errorf("handler kind %s is not a func", t.Kind())
	}

	switch t.NumOut() {
	case 2:
		if t.Out(1) == errorType {
			result.Error = true
		} else if t.Out(0) == reflect.TypeOf(http.StatusOK) {
			result.Status = true
		} else {
			return result, fmt.Errorf("return parameter type mismatch: index 0, got %s, want int [%s]", t.Out(0), dumpSignature(t))
		}
	case 3:
		result.Status = true
		result.Error = true

		if t.Out(2) != errorType {
			return result, fmt.Errorf("return parameter type mismatch: index 2, got %s, want error [%s]", t.Out(2), dumpSignature(t))
		} else if t.Out(0) != reflect.TypeOf(http.StatusOK) {
			return result, fmt.Errorf("return parameter type mismatch: index 0, got %s, want int [%s]", t.Out(0), dumpSignature(t))
		}
	case 1:
		// nop
	default:
		return result, fmt.Errorf("unknown handler signature: %s", dumpSignature(t))
	}

	if t.NumIn() == 2 {
		result.InParameters = true
		if t.In(0) != headerType {
			return result, fmt.Errorf("input parameter type mismatch: index 0, got %s, want autojson.HeaderHandler [%s]", t.In(0), dumpSignature(t))
		}

		if (t.In(1) != reflect.TypeOf(&http.Request{})) {
			return result, fmt.Errorf("input parameter type mismatch: index 1, got %s, want *http.Request [%s]", t.In(1), dumpSignature(t))
		}
	} else if t.NumIn() != 0 {
		return result, fmt.Errorf("unknown handler signature: %s", dumpSignature(t))
	}

	return result, nil
}
