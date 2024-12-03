package form

import (
	"context"
)

type Context struct {
	Query map[string]any
}

func (ctx *Context) Get(key string) string {
	if val, ok := ctx.Query[key]; ok {
		switch val := val.(type) {
		case []string:
			return val[0]
		case string:
			return val
		}
	}
	return ""
}

func UserContext(ctx context.Context) *Context {
	if userContext, ok := ctx.Value("UserContext").(*Context); ok {
		return userContext
	}
	return &Context{Query: make(map[string]any)}
}
