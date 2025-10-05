package form

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"strings"
)

const (
	MIMEApplicationJSON = "application/json"
	MIMEApplicationForm = "application/x-www-form-urlencoded"
	MIMEMultipartForm   = "multipart/form-data"
	MIMETextPlain       = "text/plain"
)

// ParseBodyAsJSON handles Fiber requests and extracts context and JSON payload.
func ParseBodyAsJSON(ctx context.Context, contentType string, body []byte, queryParams map[string]string) (context.Context, []byte, error) {
	userContext := &Context{Query: make(map[string]any)}
	for key, value := range queryParams {
		userContext.Query[key] = value
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
	case strings.Contains(contentType, MIMEMultipartForm):
		if body == nil {
			return ctx, nil, errors.New("empty multipart body")
		}
		_, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			return ctx, nil, fmt.Errorf("failed to parse content type: %v", err)
		}
		boundary, ok := params["boundary"]
		if !ok {
			return ctx, nil, errors.New("no boundary in multipart content type")
		}
		reader := multipart.NewReader(strings.NewReader(string(body)), boundary)
		form := make(map[string]any)
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return ctx, nil, fmt.Errorf("failed to read multipart part: %v", err)
			}
			name := part.FormName()
			if name == "" {
				continue // skip files
			}
			valueBytes, err := io.ReadAll(part)
			if err != nil {
				return ctx, nil, fmt.Errorf("failed to read part value: %v", err)
			}
			valueStr := string(valueBytes)
			if existing, ok := form[name]; ok {
				if sl, ok := existing.([]string); ok {
					form[name] = append(sl, valueStr)
				} else {
					form[name] = []string{existing.(string), valueStr}
				}
			} else {
				form[name] = valueStr
			}
			part.Close()
		}
		for key, v := range form {
			userContext.Query[key] = v
		}
		result = userContext.Query
	case strings.Contains(contentType, MIMETextPlain):
		if len(body) == 0 {
			return ctx, nil, nil
		}
		result = string(body)
	default:
		return ctx, body, nil
	}
	bt, err := json.Marshal(result)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, bt, nil
}
