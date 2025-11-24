package ast

import (
	"fmt"
	"sort"
)

type NodeType string

const (
	Added     NodeType = "added"
	Removed   NodeType = "removed"
	Unchanged NodeType = "unchanged"
	Updated   NodeType = "updated"
	Nested    NodeType = "nested"
)

type Node struct {
	Key      string
	Action   NodeType
	OldVal   any
	NewVal   any
	Children []Node
}

type JsonNode struct {
	Key      string     `json:"key"`
	Type     string     `json:"type"`
	OldValue any        `json:"oldValue,omitempty"`
	NewValue any        `json:"newValue,omitempty"`
	Children []JsonNode `json:"children,omitempty"`
}

func BuildDiff(a, b map[string]any) []Node {
	keys := unionKeys(a, b)
	sort.Strings(keys)
	out := make([]Node, 0, len(keys))
	for _, k := range keys {
		v1, ok1 := a[k]
		v2, ok2 := b[k]
		switch {
		case ok1 && !ok2:
			out = append(out, Node{Key: k, Action: Removed, OldVal: v1})
		case !ok1 && ok2:
			out = append(out, Node{Key: k, Action: Added, NewVal: v2})
		default:
			if m1, ok := v1.(map[string]any); ok {
				if m2, ok2 := v2.(map[string]any); ok2 {
					out = append(out, Node{Key: k, Action: Nested, Children: BuildDiff(m1, m2)})
					continue
				}
			}
			if equals(v1, v2) {
				out = append(out, Node{Key: k, Action: Unchanged, OldVal: v1})
			} else {
				out = append(out, Node{Key: k, Action: Updated, OldVal: v1, NewVal: v2})
			}
		}
	}
	return out
}

func unionKeys(a, b map[string]any) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}

func equals(a, b any) bool {
	return fmt.Sprintf("%#v", a) == fmt.Sprintf("%#v", b)
}
