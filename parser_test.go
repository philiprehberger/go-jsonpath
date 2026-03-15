package jsonpath

import (
	"testing"
)

func TestParse_Root(t *testing.T) {
	tokens, err := parse("$")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].kind != tokenRoot {
		t.Fatalf("expected root token, got %v", tokens[0].kind)
	}
}

func TestParse_SimpleField(t *testing.T) {
	tokens, err := parse("$.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[1].kind != tokenField || tokens[1].field != "name" {
		t.Fatalf("expected field 'name', got %+v", tokens[1])
	}
}

func TestParse_NestedFields(t *testing.T) {
	tokens, err := parse("$.a.b.c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 4 {
		t.Fatalf("expected 4 tokens, got %d", len(tokens))
	}
	expected := []string{"a", "b", "c"}
	for i, name := range expected {
		tok := tokens[i+1]
		if tok.kind != tokenField || tok.field != name {
			t.Fatalf("token %d: expected field '%s', got %+v", i+1, name, tok)
		}
	}
}

func TestParse_ArrayIndex(t *testing.T) {
	tokens, err := parse("$.items[0]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
	if tokens[2].kind != tokenIndex || tokens[2].index != 0 {
		t.Fatalf("expected index 0, got %+v", tokens[2])
	}
}

func TestParse_Wildcard(t *testing.T) {
	tokens, err := parse("$.items[*]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
	if tokens[2].kind != tokenWildcard {
		t.Fatalf("expected wildcard token, got %+v", tokens[2])
	}
}

func TestParse_Complex(t *testing.T) {
	tokens, err := parse("$.store.books[0].title")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 5 {
		t.Fatalf("expected 5 tokens, got %d", len(tokens))
	}
	if tokens[1].kind != tokenField || tokens[1].field != "store" {
		t.Fatalf("token 1: expected field 'store', got %+v", tokens[1])
	}
	if tokens[2].kind != tokenField || tokens[2].field != "books" {
		t.Fatalf("token 2: expected field 'books', got %+v", tokens[2])
	}
	if tokens[3].kind != tokenIndex || tokens[3].index != 0 {
		t.Fatalf("token 3: expected index 0, got %+v", tokens[3])
	}
	if tokens[4].kind != tokenField || tokens[4].field != "title" {
		t.Fatalf("token 4: expected field 'title', got %+v", tokens[4])
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []string{
		"",
		"name",
		"$.items[",
		"$.items[abc]",
		"$.",
		"$.items[-1]",
	}
	for _, path := range cases {
		_, err := parse(path)
		if err == nil {
			t.Errorf("expected error for path %q, got nil", path)
		}
	}
}
