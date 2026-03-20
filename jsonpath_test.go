package jsonpath

import (
	"encoding/json"
	"testing"
)

func TestGet_SimpleField(t *testing.T) {
	data := []byte(`{"name":"Alice"}`)
	got, err := Get[string](data, "$.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Alice" {
		t.Fatalf("expected 'Alice', got %q", got)
	}
}

func TestGet_NestedField(t *testing.T) {
	data := []byte(`{"user":{"name":"Alice"}}`)
	got, err := Get[string](data, "$.user.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Alice" {
		t.Fatalf("expected 'Alice', got %q", got)
	}
}

func TestGet_ArrayIndex(t *testing.T) {
	data := []byte(`{"users":[{"name":"Alice"},{"name":"Bob"}]}`)
	got, err := Get[string](data, "$.users[0].name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Alice" {
		t.Fatalf("expected 'Alice', got %q", got)
	}
}

func TestGet_IntValue(t *testing.T) {
	data := []byte(`{"count":42}`)
	got, err := Get[float64](data, "$.count")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestGet_BoolValue(t *testing.T) {
	data := []byte(`{"active":true}`)
	got, err := Get[bool](data, "$.active")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Fatalf("expected true, got %v", got)
	}
}

func TestGetAll_Wildcard(t *testing.T) {
	data := []byte(`{"users":[{"name":"Alice"},{"name":"Bob"},{"name":"Charlie"}]}`)
	got, err := GetAll[string](data, "$.users[*].name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"Alice", "Bob", "Charlie"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(got))
	}
	for i, v := range expected {
		if got[i] != v {
			t.Fatalf("index %d: expected %q, got %q", i, v, got[i])
		}
	}
}

func TestGetRaw(t *testing.T) {
	data := []byte(`{"name":"Alice","age":30}`)
	got, err := GetRaw(data, "$.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatalf("expected string, got %T", got)
	}
	if s != "Alice" {
		t.Fatalf("expected 'Alice', got %q", s)
	}
}

func TestSet_SimpleField(t *testing.T) {
	data := []byte(`{"name":"Alice"}`)
	result, err := Set(data, "$.name", "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(result, &m); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if m["name"] != "Bob" {
		t.Fatalf("expected 'Bob', got %v", m["name"])
	}
}

func TestSet_NestedField(t *testing.T) {
	data := []byte(`{"user":{"name":"Alice","age":30}}`)
	result, err := Set(data, "$.user.name", "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := Get[string](result, "$.user.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Bob" {
		t.Fatalf("expected 'Bob', got %q", got)
	}
}

func TestGet_PathNotFound(t *testing.T) {
	data := []byte(`{"name":"Alice"}`)
	_, err := Get[string](data, "$.missing")
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

func TestGet_InvalidPath(t *testing.T) {
	data := []byte(`{"name":"Alice"}`)
	_, err := Get[string](data, "invalid")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestExists(t *testing.T) {
	data := []byte(`{"name":"Alice","address":{"city":"NYC"}}`)

	ok, err := Exists(data, "$.name")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected name to exist")
	}

	ok, err = Exists(data, "$.missing")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected missing to not exist")
	}

	ok, err = Exists(data, "$.address.city")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected address.city to exist")
	}
}

func TestDelete(t *testing.T) {
	data := []byte(`{"name":"Alice","age":30}`)

	result, err := Delete(data, "$.age")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Get[int](result, "$.age")
	if err == nil {
		t.Error("expected age to be deleted")
	}

	name, err := Get[string](result, "$.name")
	if err != nil {
		t.Fatal(err)
	}
	if name != "Alice" {
		t.Errorf("expected Alice, got %s", name)
	}
}
