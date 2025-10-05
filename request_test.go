package form

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestParseBodyAsJSON_JSONObject(t *testing.T) {
	ctx := context.Background()
	contentType := "application/json"
	body := []byte(`{"name": "John", "age": 30}`)
	queryParams := map[string]string{"query": "test"}

	newCtx, result, err := ParseBodyAsJSON(ctx, contentType, body, queryParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	userCtx := UserContext(newCtx)
	if userCtx.Get("query") != "test" {
		t.Errorf("Expected query param 'test', got %s", userCtx.Get("query"))
	}

	var parsed map[string]any
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	expected := map[string]any{"name": "John", "age": float64(30)}
	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
}

func TestParseBodyAsJSON_JSONArray(t *testing.T) {
	ctx := context.Background()
	contentType := "application/json"
	body := []byte(`[{"name": "John"}, {"name": "Jane"}]`)
	queryParams := map[string]string{}

	_, result, err := ParseBodyAsJSON(ctx, contentType, body, queryParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var parsed []map[string]any
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	expected := []map[string]any{{"name": "John"}, {"name": "Jane"}}
	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
}

func TestParseBodyAsJSON_FormURLEncoded(t *testing.T) {
	ctx := context.Background()
	contentType := "application/x-www-form-urlencoded"
	body := []byte("name=John&age=30")
	queryParams := map[string]string{"query": "test"}

	newCtx, result, err := ParseBodyAsJSON(ctx, contentType, body, queryParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	userCtx := UserContext(newCtx)
	if userCtx.Get("name") != "John" {
		t.Errorf("Expected name 'John', got %s", userCtx.Get("name"))
	}
	if userCtx.Get("query") != "test" {
		t.Errorf("Expected query 'test', got %s", userCtx.Get("query"))
	}

	var parsed map[string]any
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	expected := map[string]any{"name": "John", "age": "30", "query": "test"}
	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
}

func TestParseBodyAsJSON_MultipartFormData(t *testing.T) {
	ctx := context.Background()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "John")
	writer.WriteField("age", "30")
	writer.Close()

	contentType := writer.FormDataContentType()
	queryParams := map[string]string{}

	newCtx, result, err := ParseBodyAsJSON(ctx, contentType, body.Bytes(), queryParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	userCtx := UserContext(newCtx)
	if userCtx.Get("name") != "John" {
		t.Errorf("Expected name 'John', got %s", userCtx.Get("name"))
	}

	var parsed map[string]any
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	expected := map[string]any{"name": "John", "age": "30"}
	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
}

func TestParseBodyAsJSON_TextPlain(t *testing.T) {
	ctx := context.Background()
	contentType := "text/plain"
	body := []byte("Hello World")
	queryParams := map[string]string{}

	_, result, err := ParseBodyAsJSON(ctx, contentType, body, queryParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var parsed string
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if parsed != "Hello World" {
		t.Errorf("Expected 'Hello World', got %s", parsed)
	}
}

func TestParseBodyAsJSON_EmptyBody(t *testing.T) {
	ctx := context.Background()
	contentType := "application/json"
	body := []byte{}
	queryParams := map[string]string{}

	_, result, err := ParseBodyAsJSON(ctx, contentType, body, queryParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result for empty body, got %v", result)
	}
}

func TestParseBodyAsJSON_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	contentType := "application/json"
	body := []byte(`{"invalid": json}`)
	queryParams := map[string]string{}

	_, _, err := ParseBodyAsJSON(ctx, contentType, body, queryParams)
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

func TestParseBodyAsJSON_UnsupportedContentType(t *testing.T) {
	ctx := context.Background()
	contentType := "application/xml"
	body := []byte("<xml></xml>")
	queryParams := map[string]string{}

	_, result, err := ParseBodyAsJSON(ctx, contentType, body, queryParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(result, body) {
		t.Errorf("Expected raw body for unsupported content type")
	}
}
