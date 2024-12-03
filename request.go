package form

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	MIMEApplicationJSON = "application/json"
	MIMEApplicationForm = "application/x-www-form-urlencoded"
)

// ParseBodyAsJSON handles Fiber requests and extracts context and JSON payload.
func ParseBodyAsJSON(ctx context.Context, contentType string, body []byte, queryParams map[string]string) (context.Context, []byte, error) {
	userContext := &Context{Query: make(map[string]any)}
	if queryParams != nil {
		for key, value := range queryParams {
			userContext.Query[key] = value
		}
	}
	ctx = context.WithValue(ctx, "UserContext", userContext)
	var result any
	switch {
	case strings.Contains(contentType, MIMEApplicationJSON):
		if len(body) == 0 {
			return ctx, nil, nil
		}
		var temp any
		if err := json.Unmarshal(body, &temp); err != nil {
			return ctx, nil, fmt.Errorf("failed to parse body: %v", err)
		}
		switch v := temp.(type) {
		case map[string]any:
			result = v
		case []any:
			parsedArray := make([]map[string]any, len(v))
			for i, item := range v {
				obj, ok := item.(map[string]any)
				if !ok {
					return ctx, nil, fmt.Errorf("invalid JSON array item at index %d", i)
				}
				parsedArray[i] = obj
			}
			result = parsedArray
		default:
			return ctx, nil, fmt.Errorf("unsupported JSON structure: %T", v)
		}
	case strings.Contains(contentType, MIMEApplicationForm):
		if body == nil {
			return ctx, nil, errors.New("empty form body")
		}
		val, err := DecodeForm(body)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to parse form data: %v", err.Error())
		}
		for key, v := range val {
			userContext.Query[key] = v
		}
		result = userContext.Query
	default:
		return ctx, body, nil
	}
	bt, err := json.Marshal(result)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, bt, nil
}
