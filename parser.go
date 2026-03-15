package jsonpath

import (
	"fmt"
	"strconv"
	"strings"
)

// tokenKind represents the type of a path token.
type tokenKind int

const (
	tokenRoot     tokenKind = iota // $
	tokenField                     // .name
	tokenIndex                     // [0]
	tokenWildcard                  // [*]
)

// token represents a single element of a parsed JSONPath expression.
type token struct {
	kind  tokenKind
	field string // field name for tokenField
	index int    // array index for tokenIndex
}

// parse converts a JSONPath expression string into a slice of tokens.
// Supported syntax: $, $.field, $.field.nested, $.array[0], $.array[*], $[0].field
func parse(path string) ([]token, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	if path[0] != '$' {
		return nil, fmt.Errorf("path must start with $: %s", path)
	}

	tokens := []token{{kind: tokenRoot}}
	remaining := path[1:]

	for len(remaining) > 0 {
		switch remaining[0] {
		case '.':
			remaining = remaining[1:]
			if len(remaining) == 0 {
				return nil, fmt.Errorf("unexpected end of path after '.': %s", path)
			}
			// Read field name until next '.', '[', or end
			end := strings.IndexAny(remaining, ".[")
			if end == -1 {
				end = len(remaining)
			}
			field := remaining[:end]
			if field == "" {
				return nil, fmt.Errorf("empty field name in path: %s", path)
			}
			tokens = append(tokens, token{kind: tokenField, field: field})
			remaining = remaining[end:]

		case '[':
			remaining = remaining[1:]
			closeBracket := strings.IndexByte(remaining, ']')
			if closeBracket == -1 {
				return nil, fmt.Errorf("unclosed bracket in path: %s", path)
			}
			content := remaining[:closeBracket]
			remaining = remaining[closeBracket+1:]

			if content == "*" {
				tokens = append(tokens, token{kind: tokenWildcard})
			} else {
				idx, err := strconv.Atoi(content)
				if err != nil {
					return nil, fmt.Errorf("invalid array index '%s' in path: %s", content, path)
				}
				if idx < 0 {
					return nil, fmt.Errorf("negative array index %d in path: %s", idx, path)
				}
				tokens = append(tokens, token{kind: tokenIndex, index: idx})
			}

		default:
			return nil, fmt.Errorf("unexpected character '%c' in path: %s", remaining[0], path)
		}
	}

	return tokens, nil
}
