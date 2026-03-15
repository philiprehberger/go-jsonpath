// Package jsonpath provides JSONPath query and extraction for Go.
// It supports type-safe extraction with generics, array indexing,
// wildcards, and setting values at paths. Zero external dependencies.
package jsonpath

import (
	"encoding/json"
	"fmt"
)

// Get extracts a single value at the given JSONPath and unmarshals it to type T.
// Returns an error if the path is invalid, not found, or the value cannot be
// converted to T.
func Get[T any](data []byte, path string) (T, error) {
	var zero T

	raw, err := GetRaw(data, path)
	if err != nil {
		return zero, err
	}

	return convert[T](raw)
}

// GetRaw extracts a value at the given JSONPath without type conversion.
// The returned value is one of: nil, bool, float64, string, []any, or map[string]any.
func GetRaw(data []byte, path string) (any, error) {
	tokens, err := parse(path)
	if err != nil {
		return nil, err
	}

	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return traverse(root, tokens, path)
}

// GetAll extracts all matching values for a wildcard JSONPath and unmarshals
// each to type T. Use this with paths containing [*].
func GetAll[T any](data []byte, path string) ([]T, error) {
	tokens, err := parse(path)
	if err != nil {
		return nil, err
	}

	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	results, err := traverseAll(root, tokens, path)
	if err != nil {
		return nil, err
	}

	out := make([]T, 0, len(results))
	for _, r := range results {
		v, err := convert[T](r)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}

	return out, nil
}

// Set sets a value at the given JSONPath and returns the modified JSON.
// The path must point to an existing location or a direct child of an
// existing object/array.
func Set(data []byte, path string, value any) ([]byte, error) {
	tokens, err := parse(path)
	if err != nil {
		return nil, err
	}

	if len(tokens) < 2 {
		return nil, fmt.Errorf("path must specify a field or index to set: %s", path)
	}

	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if err := setValue(root, tokens[1:], value, path); err != nil {
		return nil, err
	}

	return json.Marshal(root)
}

// traverse walks the data structure according to the token path and returns
// the value at the final token.
func traverse(data any, tokens []token, fullPath string) (any, error) {
	current := data

	for i, tok := range tokens {
		switch tok.kind {
		case tokenRoot:
			continue

		case tokenField:
			obj, ok := current.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("path not found: %s (expected object at token %d)", fullPath, i)
			}
			val, exists := obj[tok.field]
			if !exists {
				return nil, fmt.Errorf("path not found: %s (field '%s' does not exist)", fullPath, tok.field)
			}
			current = val

		case tokenIndex:
			arr, ok := current.([]any)
			if !ok {
				return nil, fmt.Errorf("path not found: %s (expected array at token %d)", fullPath, i)
			}
			if tok.index >= len(arr) {
				return nil, fmt.Errorf("path not found: %s (index %d out of range, length %d)", fullPath, tok.index, len(arr))
			}
			current = arr[tok.index]

		case tokenWildcard:
			return nil, fmt.Errorf("use GetAll for wildcard paths: %s", fullPath)
		}
	}

	return current, nil
}

// traverseAll walks the data structure and collects all matches for wildcard paths.
func traverseAll(data any, tokens []token, fullPath string) ([]any, error) {
	current := []any{data}

	for i, tok := range tokens {
		var next []any

		switch tok.kind {
		case tokenRoot:
			continue

		case tokenField:
			for _, c := range current {
				obj, ok := c.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("path not found: %s (expected object at token %d)", fullPath, i)
				}
				val, exists := obj[tok.field]
				if !exists {
					return nil, fmt.Errorf("path not found: %s (field '%s' does not exist)", fullPath, tok.field)
				}
				next = append(next, val)
			}

		case tokenIndex:
			for _, c := range current {
				arr, ok := c.([]any)
				if !ok {
					return nil, fmt.Errorf("path not found: %s (expected array at token %d)", fullPath, i)
				}
				if tok.index >= len(arr) {
					return nil, fmt.Errorf("path not found: %s (index %d out of range)", fullPath, tok.index)
				}
				next = append(next, arr[tok.index])
			}

		case tokenWildcard:
			for _, c := range current {
				arr, ok := c.([]any)
				if !ok {
					return nil, fmt.Errorf("path not found: %s (expected array at token %d)", fullPath, i)
				}
				next = append(next, arr...)
			}
		}

		current = next
	}

	return current, nil
}

// setValue sets a value at the path specified by tokens within the data structure.
func setValue(data any, tokens []token, value any, fullPath string) error {
	if len(tokens) == 0 {
		return fmt.Errorf("cannot set root value: %s", fullPath)
	}

	// Navigate to the parent of the target
	current := data
	for _, tok := range tokens[:len(tokens)-1] {
		switch tok.kind {
		case tokenField:
			obj, ok := current.(map[string]any)
			if !ok {
				return fmt.Errorf("path not found: %s (expected object)", fullPath)
			}
			val, exists := obj[tok.field]
			if !exists {
				return fmt.Errorf("path not found: %s (field '%s' does not exist)", fullPath, tok.field)
			}
			current = val

		case tokenIndex:
			arr, ok := current.([]any)
			if !ok {
				return fmt.Errorf("path not found: %s (expected array)", fullPath)
			}
			if tok.index >= len(arr) {
				return fmt.Errorf("path not found: %s (index %d out of range)", fullPath, tok.index)
			}
			current = arr[tok.index]

		default:
			return fmt.Errorf("cannot set value through wildcard path: %s", fullPath)
		}
	}

	// Set the value at the final token
	last := tokens[len(tokens)-1]
	switch last.kind {
	case tokenField:
		obj, ok := current.(map[string]any)
		if !ok {
			return fmt.Errorf("path not found: %s (expected object)", fullPath)
		}
		obj[last.field] = value

	case tokenIndex:
		arr, ok := current.([]any)
		if !ok {
			return fmt.Errorf("path not found: %s (expected array)", fullPath)
		}
		if last.index >= len(arr) {
			return fmt.Errorf("path not found: %s (index %d out of range)", fullPath, last.index)
		}
		arr[last.index] = value

	default:
		return fmt.Errorf("cannot set value at wildcard: %s", fullPath)
	}

	return nil
}

// convert marshals a value to JSON and unmarshals it into type T.
func convert[T any](value any) (T, error) {
	var zero T

	b, err := json.Marshal(value)
	if err != nil {
		return zero, fmt.Errorf("failed to convert value: %w", err)
	}

	var result T
	if err := json.Unmarshal(b, &result); err != nil {
		return zero, fmt.Errorf("failed to convert value to target type: %w", err)
	}

	return result, nil
}
