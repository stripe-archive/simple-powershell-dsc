package regexpat

import (
	"context"
	"net/http"

	"goji.io/pattern"
)

type regexMatch struct {
	context.Context
	matches map[pattern.Variable]string
}

func (m *regexMatch) Value(key interface{}) interface{} {
	if key == pattern.AllVariables {
		return m.matches
	}

	v, ok := key.(pattern.Variable)
	if !ok {
		return m.Context.Value(key)
	}

	return m.matches[v]
}

func Param(r *http.Request, name string) string {
	return r.Context().Value(pattern.Variable(name)).(string)
}
